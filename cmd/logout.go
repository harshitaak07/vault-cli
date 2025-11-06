package cmd

import (
	"fmt"

	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "End the current session",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = session.Clear()
		fmt.Println("Logged out.")
		return nil
	},
}
