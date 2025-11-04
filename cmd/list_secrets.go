package cmd

import (
	"fmt"

	"vault-cli/internal/secrets"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var cat string

var listSecretsCmd = &cobra.Command{
	Use:   "list-secrets",
	Short: "List stored secrets (optionally filter by category)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := session.Require(); err != nil {
			return err
		}
		items, err := secrets.List(database, cat)
		if err != nil {
			return err
		}
		for _, s := range items {
			fmt.Printf("- %s / %s\n", s.Category, s.Name)
		}
		return nil
	},
}

func init() {
	listSecretsCmd.Flags().StringVar(&cat, "cat", "", "Filter by category")
}
