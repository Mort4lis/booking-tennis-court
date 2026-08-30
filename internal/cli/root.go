package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd builds the CLI command tree. version is injected from main via
// -ldflags at build time. Each subcommand loads exactly the configuration it
// needs, so there is no shared PersistentPreRunE.
func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "bot",
		Short:         "Telegram bot for booking a tennis court",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newServeCmd(),
		newMigrateCmd(),
		newVersionCmd(version),
	)

	return root
}
