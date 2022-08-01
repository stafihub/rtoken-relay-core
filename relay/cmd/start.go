package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cosmosChain "github.com/stafihub/cosmos-relay-sdk/chain"
	"github.com/stafihub/rtoken-relay-core/common/config"
	"github.com/stafihub/rtoken-relay-core/common/core"
	"github.com/stafihub/rtoken-relay-core/common/log"
	stafiHubChain "github.com/stafihub/stafi-hub-relay-sdk/chain"
	stafiHubXLedgerTypes "github.com/stafihub/stafihub/x/ledger/types"
)

const (
	flagConfig   = "config"
	flagLogLevel = "log_level"
)

var defaultConfigPath = os.ExpandEnv("./config.json")

func startCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start relay procedure",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}
			fmt.Printf("config path: %s\n", configPath)
			logLevelStr, err := cmd.Flags().GetString(flagLogLevel)
			if err != nil {
				return err
			}
			logLevel, err := logrus.ParseLevel(logLevelStr)
			if err != nil {
				return err
			}
			logrus.SetLevel(logLevel)

			cfg, err := config.GetConfig(configPath)
			if err != nil {
				return err
			}

			err = log.InitLogFile(cfg.LogFilePath)
			if err != nil {
				return err
			}

			// Used to signal core shutdown due to fatal error
			sysErr := make(chan error)
			c := core.NewCore(log.NewLog(), sysErr)

			// ======================== init stafiHub
			stafiHubChainConfig := cfg.NativeChain
			stafiHubChainConfig.Rsymbol = string(core.HubRFIS)
			bts, err := json.Marshal(stafiHubChainConfig.Opts)
			if err != nil {
				return err
			}
			option := stafiHubChain.ConfigOption{}
			err = json.Unmarshal(bts, &option)
			if err != nil {
				return err
			}
			option.CaredSymbol = cfg.ExternalChain.Rsymbol
			option.BlockstorePath = cfg.BlockstorePath

			stafiHubChainConfig.Opts = option
			stafiHubChain := stafiHubChain.NewChain()
			err = stafiHubChain.Initialize(&stafiHubChainConfig, log.NewLog("chain", stafiHubChainConfig.Name), sysErr)
			if err != nil {
				return err
			}
			c.AddChain(stafiHubChain)

			//========================== init external chain
			chainConfig := cfg.ExternalChain
			var newChain core.Chain

			// load option config from file
			bts, err = json.Marshal(chainConfig.Opts)
			if err != nil {
				return err
			}
			cosmosOption := cosmosChain.ConfigOption{}
			err = json.Unmarshal(bts, &cosmosOption)
			if err != nil {
				return err
			}

			cosmosOption.BlockstorePath = cfg.BlockstorePath

			// prepare r params from stafihub
			rParams, err := stafiHubChain.GetRParams(chainConfig.Rsymbol)
			if err != nil {
				return err
			}
			poolRes, err := stafiHubChain.GetPools(rParams.RParams.Denom)
			if err != nil {
				return err
			}

			icaPoolsRes, err := stafiHubChain.GetIcaPools(rParams.RParams.Denom)
			if err != nil {
				return err
			}
			if len(cosmosOption.PoolNameSubKey) == 0 && len(icaPoolsRes.IcaPoolList) == 0 {
				return fmt.Errorf("no pool")
			}

			icaPoolsMap := make(map[string]string)
			icaPoolCtrlChannel := make(map[string]string)

			for _, value := range icaPoolsRes.IcaPoolList {
				if value.Status == stafiHubXLedgerTypes.IcaPoolStatusSetWithdrawal {
					icaPoolsMap[value.DelegationAccount.Address] = value.WithdrawalAccount.Address
					icaPoolCtrlChannel[value.DelegationAccount.Address] = value.DelegationAccount.CtrlChannelId
				}
			}

			cosmosOption.PoolAddressThreshold = make(map[string]uint32)
			cosmosOption.PoolTargetValidators = make(map[string][]string)
			bondedIcaPoolWithdrawalAddr := make(map[string]string)

			for _, poolAddressStr := range poolRes.GetAddrs() {
				// filter icapool
				if withdrawalAddr, exist := icaPoolsMap[poolAddressStr]; exist {
					bondedIcaPoolWithdrawalAddr[poolAddressStr] = withdrawalAddr
				} else {
					// get pool threshold
					poolDetail, err := stafiHubChain.GetPoolDetail(rParams.RParams.Denom, poolAddressStr)
					if err != nil {
						return err
					}
					if poolDetail.Detail.Threshold <= 0 {
						return fmt.Errorf("pool threshold is zero in stafihub, pool: %s", poolAddressStr)
					}
					cosmosOption.PoolAddressThreshold[poolAddressStr] = poolDetail.Detail.Threshold
				}

				// get pool targetValidators from rvalidator
				selectedValidators, err := stafiHubChain.GetSelectedValidators(rParams.RParams.Denom, poolAddressStr)
				if err != nil {
					return err
				}
				if len(selectedValidators.RValidatorList) <= 0 {
					return fmt.Errorf("pool selected validators is empty, pool: %s", poolAddressStr)
				}
				cosmosOption.PoolTargetValidators[poolAddressStr] = selectedValidators.RValidatorList
			}

			cosmosOption.EraSeconds = rParams.RParams.EraSeconds
			cosmosOption.GasPrice = rParams.RParams.GasPrice
			cosmosOption.LeastBond = rParams.RParams.LeastBond
			cosmosOption.Offset = rParams.RParams.Offset

			cosmosOption.IcaPoolWithdrawalAddr = bondedIcaPoolWithdrawalAddr
			cosmosOption.IcaPoolCtrlChannel = icaPoolCtrlChannel

			// prepare account prefix from stafihub
			prefixRes, err := stafiHubChain.GetAddressPrefix(chainConfig.Rsymbol)
			if err != nil {
				return err
			}
			cosmosOption.AccountPrefix = prefixRes.GetAccAddressPrefix()

			chainConfig.Opts = cosmosOption
			newChain = cosmosChain.NewChain()
			err = newChain.Initialize(&chainConfig, log.NewLog("chain", chainConfig.Name), sysErr)
			if err != nil {
				return fmt.Errorf("newChain.Initialize failed: %s", err)
			}
			c.AddChain(newChain)

			// =============== start
			c.Start()

			return nil
		},
	}

	cmd.Flags().String(flagConfig, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}
