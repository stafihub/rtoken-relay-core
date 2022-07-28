package cmd

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
)

var (
	defaultNodeHome = "./keys/stafihub"

	flagPrefix = "prefix"
	flagHome   = "home"
)

func keyCmd() *cobra.Command {
	encodingConfig := MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin)

	keysCmd := &cobra.Command{
		Use:   "keys",
		Short: "Key tool to manage keys",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			prefix, err := cmd.Flags().GetString(flagPrefix)
			if err != nil {
				return err
			}

			SetPrefixes(prefix)

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return fmt.Errorf("SetCmdClientContextHandler err: %s", err)
			}
			return nil
		},
	}

	keysCmd.AddCommand(
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

	keysCmd.PersistentFlags().String(flagPrefix, "stafi", "The chain prefix")
	keysCmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	keysCmd.PersistentFlags().String(flags.FlagKeyringBackend, "file", "Select keyring's backend (os|file|test)")
	keysCmd.PersistentFlags().String("output", "text", "Output format (text|json)")
	keysCmd.PersistentFlags().StringP(flagHome, "", defaultNodeHome, "Directory for config and data")
	return keysCmd
}
