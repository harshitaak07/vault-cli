package cmd

import (
	"fmt"

	"vault-cli/internal/secrets"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var rotateKeysCmd = &cobra.Command{
	Use:   "rotate-keys",
	Short: "Re-encrypt all secrets with fresh data keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := session.Require(); err != nil {
			return err
		}
		n, err := secrets.Rotate(database, cfg)
		if err != nil {
			return err
		}
		fmt.Printf("🔁 Rotated %d secrets.\n", n)
		return nil
	},
}
