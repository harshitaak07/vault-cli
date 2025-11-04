package cmd

import (
	"fmt"

	"vault-cli/internal/secrets"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var getSecretCmd = &cobra.Command{
	Use:   "get-secret <category> <name>",
	Short: "Retrieve and decrypt a secret",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := session.Require(); err != nil {
			return err
		}
		val, err := secrets.Get(database, cfg, args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Printf("%s = %s\n", args[1], val)
		return nil
	},
}
