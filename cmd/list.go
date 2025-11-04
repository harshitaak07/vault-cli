package cmd

import (
	"log"

	"vault-cli/internal/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored file metadata from local database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.PrintDBEntries(database); err != nil {
			log.Fatalf("list: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
