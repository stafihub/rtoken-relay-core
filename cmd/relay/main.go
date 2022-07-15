package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	cosmosChain "github.com/stafihub/cosmos-relay-sdk/chain"
	"github.com/stafihub/rtoken-relay-core/common/config"
	"github.com/stafihub/rtoken-relay-core/common/core"
	"github.com/stafihub/rtoken-relay-core/common/log"
	stafiHubChain "github.com/stafihub/stafi-hub-relay-sdk/chain"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var mainFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2022 Stafi Protocol Authors"
	app.Name = "relay"
	app.Usage = "relay"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2022"}}
	app.Version = "2.0.0"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{}

	app.Flags = append(app.Flags, mainFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}
	verbosity := ctx.String(config.VerbosityFlag.Name)
	if verbosity != "" {
		level, err := logrus.ParseLevel(verbosity)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
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
	if len(cosmosOption.PoolNameSubKey) == 0 {
		return fmt.Errorf("no pool and subkey")
	}

	// prepare r params from stafihub
	rParams, err := stafiHubChain.GetRParams(chainConfig.Rsymbol)
	if err != nil {
		return err
	}
	poolRes, err := stafiHubChain.GetPools(rParams.RParams.Denom)
	if err != nil {
		return err
	}

	cosmosOption.PoolAddressThreshold = make(map[string]uint32)
	cosmosOption.PoolTargetValidators = make(map[string][]string)
	for _, poolAddressStr := range poolRes.GetAddrs() {
		// get pool threshold
		poolDetail, err := stafiHubChain.GetPoolDetail(rParams.RParams.Denom, poolAddressStr)
		if err != nil {
			return err
		}
		if poolDetail.Detail.Threshold <= 0 {
			return fmt.Errorf("pool threshold is zero in stafihub, pool: %s", poolAddressStr)
		}
		cosmosOption.PoolAddressThreshold[poolAddressStr] = poolDetail.Detail.Threshold

		// get pool targetValidators
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
}
