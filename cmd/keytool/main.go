package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"os"
)

const DefaultHomeDir = "./keys/safi_hub"

func main() {
	encodingConfig := MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(xAuthTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(DefaultHomeDir).WithKeyringOptions()

	rootCmd := &cobra.Command{
		Use:   "keys tool",
		Short: "tool to manage keys",
	}
	genStafiCmd := &cobra.Command{
		Use:   "genstafi",
		Short: "tool to manage stafi-hub keys",
		Long:  "Notice: The keyring supports os|file|test backends, but relay now only support the file backend",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			SetPrefixes("fis")
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}
	genStafiCmd.AddCommand(keys.Commands(DefaultHomeDir))

	genCosmosCmd := &cobra.Command{
		Use:   "gencosmos",
		Short: "tool to manage cosmos-hub keys",
		Long:  "Notice: The keyring supports os|file|test backends, but relay now only support the file backend",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			SetPrefixes("cosmos")
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}
	genCosmosCmd.AddCommand(keys.Commands(DefaultHomeDir))

	rootCmd.AddCommand(genStafiCmd)
	rootCmd.AddCommand(genCosmosCmd)

	svrcmd.Execute(rootCmd, DefaultHomeDir)
}
