package cmd

import (
	"fmt"
	"time"

	"vault-cli/internal/auth"
	"vault-cli/internal/session"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Start a session (verifies master password)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := auth.VerifyPassword(cfg.PasswordFile); !ok {
			return fmt.Errorf("access denied: wrong password")
		}
		if err := session.Save("admin", 15*time.Minute); err != nil {
			return err
		}
		fmt.Println("✅ Session started (15m).")
		return nil
	},
}
