package main

import (
	"database/sql"
	"fmt"
)

func UploadHandler(file string, cfg *Config, db *sql.DB) error {
	fmt.Printf("Uploading %s using mode=%s\n", file, cfg.Mode)
	if cfg.Mode == "local" {
		return LocalUpload(file)
	}
	err := EncryptAndUpload(file, cfg, db)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	_ = LogToCloudWatch(fmt.Sprintf("File uploaded: %s", file))
	return nil
}

func DownloadHandler(file string, cfg *Config, db *sql.DB) error {
	fmt.Printf("Downloading %s using mode=%s\n", file, cfg.Mode)
	if cfg.Mode == "local" {
		return LocalDownload(file)
	}
	return DownloadAndDecrypt(file, cfg, db)
}
