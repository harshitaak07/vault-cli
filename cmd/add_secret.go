package cmd

import (
	"fmt"

	"vault-cli/internal/secrets"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var addSecretCmd = &cobra.Command{
	Use:   "add-secret <category> <name> <value>",
	Short: "Add or update an encrypted secret",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := session.Require(); err != nil {
			return err
		}
		s := secrets.Secret{Category: args[0], Name: args[1], Value: args[2]}
		if err := secrets.Add(database, cfg, s); err != nil {
			return fmt.Errorf("add-secret: %w", err)
		}
		fmt.Println("Secret stored.")
		return nil
	},
}
