package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

const StafiHubDefaultHomeDir = "./keys/stafihub"
const CosmosHubDefaultHomeDir = "./keys/cosmoshub"

func main() {
	encodingConfig := MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(xAuthTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock)

	rootCmd := &cobra.Command{
		Use:   "keytool",
		Short: "tool to manage keys",
	}
	genStafiCmd := &cobra.Command{
		Use:   "genstafi",
		Short: "tool to manage stafi-hub keys",
		Long:  "Notice: The keyring supports os|file|test backends, but relay now only support the file backend",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx = initClientCtx.WithHomeDir(StafiHubDefaultHomeDir).WithKeyringOptions()
			SetPrefixes("fis")
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}
	genStafiCmd.AddCommand(keys.Commands(StafiHubDefaultHomeDir))

	genCosmosCmd := &cobra.Command{
		Use:   "gencosmos",
		Short: "tool to manage cosmos-hub keys",
		Long:  "Notice: The keyring supports os|file|test backends, but relay now only support the file backend",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx = initClientCtx.WithHomeDir(CosmosHubDefaultHomeDir).WithKeyringOptions()
			SetPrefixes("cosmos")
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}
	genCosmosCmd.AddCommand(keys.Commands(CosmosHubDefaultHomeDir))

	rootCmd.AddCommand(genStafiCmd)
	rootCmd.AddCommand(genCosmosCmd)

	svrcmd.Execute(rootCmd, "")
}
