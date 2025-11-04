package cmd

import (
	"fmt"
	"log"

	"vault-cli/internal/local"

	"github.com/spf13/cobra"
)

var localDownloadCmd = &cobra.Command{
	Use:   "local-download <file>",
	Short: "Retrieve a file from local vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := local.LocalDownload(file); err != nil {
			log.Fatalf("local download failed: %v", err)
		}
		fmt.Println("✅ Local download successful.")
	},
}

func init() {
	rootCmd.AddCommand(localDownloadCmd)
}
