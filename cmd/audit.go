package cmd

import (
	"log"

	"vault-cli/internal/db"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show recent upload/download audit logs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.PrintAudit(database); err != nil {
			log.Fatalf("audit: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
