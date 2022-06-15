package main

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
)

var defaultNodeHome = "./keys/stafihub"

func main() {
	encodingConfig := MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin)

	rootCmd := &cobra.Command{
		Use:   "keytool",
		Short: "tool to manage keys",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				return err
			}

			SetPrefixes(prefix)
			return client.SetCmdClientContextHandler(initClientCtx, cmd)
		},
	}

	rootCmd.AddCommand(
		keys.MnemonicKeyCommand(),
		keys.AddKeyCommand(),
		keys.ExportKeyCommand(),
		keys.ImportKeyCommand(),
		keys.ListKeysCmd(),
		keys.ShowKeysCmd(),
		keys.DeleteKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
	)

	rootCmd.PersistentFlags().String("prefix", "stafi", "The chain prefix")
	rootCmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	rootCmd.PersistentFlags().String(flags.FlagKeyringBackend, "file", "Select keyring's backend (os|file|test)")
	rootCmd.PersistentFlags().String("output", "text", "Output format (text|json)")
	rootCmd.PersistentFlags().StringP("home", "", defaultNodeHome, "Directory for config and data")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
