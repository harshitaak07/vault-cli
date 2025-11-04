package cmd

import (
	"vault-cli/internal/core"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Show vault summary report (file count, total size, recent uploads)",
	Run: func(cmd *cobra.Command, args []string) {
		core.GenerateReport(database)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
