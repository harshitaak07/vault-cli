package cmd

import (
	"fmt"
	"log"

	"vault-cli/internal/local"

	"github.com/spf13/cobra"
)

var localUploadCmd = &cobra.Command{
	Use:   "local-upload <file>",
	Short: "Copy file to local vault (no cloud upload)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := local.LocalUpload(file); err != nil {
			log.Fatalf("local upload failed: %v", err)
		}
		fmt.Println("✅ Local upload successful.")
	},
}

func init() {
	rootCmd.AddCommand(localUploadCmd)
}
