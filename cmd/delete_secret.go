package cmd

import (
	"vault-cli/internal/secrets"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var deleteSecretCmd = &cobra.Command{
	Use:   "delete-secret <category> <name>",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := session.Require(); err != nil {
			return err
		}
		return secrets.Delete(database, args[0], args[1])
	},
}
