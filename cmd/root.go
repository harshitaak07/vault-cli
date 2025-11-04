package cmd

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"vault-cli/internal/auth"
	"vault-cli/internal/config"
	"vault-cli/internal/db"
)

var (
	cfg      *config.Config
	database *sql.DB
	rootCmd  = &cobra.Command{
		Use:   "vault",
		Short: "Vault CLI — Secure file encryption and storage manager",
		Long: `Vault CLI allows you to encrypt, upload, download, and manage files securely
using AWS KMS, S3, and DynamoDB or local vault mode.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_ = godotenv.Load(".env")

			var err error
			cfg, err = config.LoadConfig()
			if err != nil {
				return err
			}

			database, err = db.OpenDB(cfg.DBPath)
			if err != nil {
				return err
			}
			if err := db.InitDB(database); err != nil {
				return err
			}

			if cfg.RequirePassword {
				if ok := auth.VerifyPassword(cfg.PasswordFile); !ok {
					log.Fatal("Access denied: wrong password")
				}
			}

			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func CloseDB() {
	if database != nil {
		_ = database.Close()
	}
}
