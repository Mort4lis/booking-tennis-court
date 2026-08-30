package cli

import (
	"github.com/spf13/cobra"

	"github.com/Mort4lis/booking-tennis-court/internal/app"
)

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the Telegram bot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return app.Run(cmd.Context())
		},
	}
}
