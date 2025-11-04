package core

import (
	"database/sql"
	"fmt"

	"vault-cli/internal/aws"
	"vault-cli/internal/config"
	"vault-cli/internal/local"
)

func UploadHandler(file string, cfg *config.Config, db *sql.DB) error {
	fmt.Printf("Uploading %s using mode=%s\n", file, cfg.Mode)
	if cfg.Mode == "local" {
		return local.LocalUpload(file)
	}
	err := aws.EncryptAndUpload(file, cfg, db)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	_ = aws.LogToCloudWatch(fmt.Sprintf("File uploaded: %s", file))
	return nil
}

func DownloadHandler(file string, cfg *config.Config, db *sql.DB) error {
	fmt.Printf("Downloading %s using mode=%s\n", file, cfg.Mode)
	if cfg.Mode == "local" {
		return local.LocalDownload(file)
	}
	return aws.DownloadAndDecrypt(file, cfg, db)
}
