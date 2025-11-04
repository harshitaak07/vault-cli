package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"vault-cli/internal/core"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Encrypt and upload a file to S3 or local vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := core.UploadHandler(file, cfg, database); err != nil {
			log.Fatalf("upload failed: %v", err)
		}
		fmt.Println("✅ File uploaded successfully.")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
