package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stafihub/cosmos-relay-sdk/client"
	"github.com/stafihub/rtoken-relay-core/common/log"
)

func multisigTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisig-transfer",
		Short: "Tranfer token from multisig account",
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

			config := Config{}
			err = loadConfig(configPath, &config)
			if err != nil {
				return err
			}
			fmt.Printf("config: %+v\n\n", config)
			fmt.Printf("Will open wallet from <%s>. \nPlease ", config.KeystorePath)
			key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, config.KeystorePath, os.Stdin)
			if err != nil {
				return err
			}

			cosmosClient, err := client.NewClient(key, config.MultisigAccountName, config.GasPrice, config.Prefix, []string{config.Endpoint}, log.NewLog("client"))
			if err != nil {
				return err
			}
			toAddress, err := types.AccAddressFromBech32(config.ToAddress)
			if err != nil {
				return err
			}
			amount, err := types.ParseCoinNormalized(config.Amount)
			if err != nil {
				return err
			}

			rawTx, err := cosmosClient.GenMultiSigRawTransferTx(toAddress, types.NewCoins(amount))
			if err != nil {
				return err
			}

			account, err := cosmosClient.QueryAccount(cosmosClient.GetFromAddress())
			if err != nil {
				return err
			}

			sigs := make([][]byte, len(config.SubAccountNameList))
			for i, subKey := range config.SubAccountNameList {
				sig, err := cosmosClient.SignMultiSigRawTxWithSeq(account.GetSequence(), rawTx, subKey)
				if err != nil {
					return err
				}
				sigs[i] = sig
			}

			_, tx, err := cosmosClient.AssembleMultiSigTx(rawTx, sigs, uint32(config.Threshold))
			if err != nil {
				return err
			}
			hash, err := cosmosClient.BroadcastTx(tx)
			if err != nil {
				return err
			}
			fmt.Println("hash ", hash)
			return nil
		},
	}

	cmd.Flags().String(flagConfig, defaultConfigPath, "Config file path")
	cmd.Flags().String(flagLogLevel, logrus.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")

	return cmd
}

type Config struct {
	KeystorePath        string   `json:"keystorePath"`
	MultisigAccountName string   `json:"MultisigAccountName"`
	SubAccountNameList  []string `json:"subAccountNameList"`
	ToAddress           string   `json:"toAddress"`
	Endpoint            string   `json:"endpoint"`
	Prefix              string   `json:"prefix"`
	GasPrice            string   `json:"gasPrice"`
	Amount              string   `json:"amount"`
	Threshold           int64    `json:"threshold"`
}

func loadConfig(file string, config *Config) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if ext != ".json" {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return json.NewDecoder(f).Decode(&config)
}
