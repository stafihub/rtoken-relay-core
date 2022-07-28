package cmd

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

var (
	appName = "relay"
)

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "relay",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {

		return nil
	}

	rootCmd.AddCommand(
		startCmd(),
		versionCmd(),
		keyCmd(),
	)
	return rootCmd
}

func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// keys cmd need set clientContext pointer
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
