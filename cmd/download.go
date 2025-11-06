package cmd

import (
	"fmt"
	"log"

	"vault-cli/internal/core"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <file>",
	Short: "Download and decrypt a file from AWS S3 or the local vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := core.DownloadHandler(file, cfg, database); err != nil {
			log.Fatalf("Download failed: %v", err)
		}

		fmt.Println("File downloaded successfully.")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
