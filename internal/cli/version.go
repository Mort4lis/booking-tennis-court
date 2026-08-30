package cli

import "github.com/spf13/cobra"

// newVersionCmd adds a `version` subcommand for scripts that prefer it over the
// `--version` flag NewRootCmd already sets.
func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the binary version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println(version)

			return nil
		},
	}
}
