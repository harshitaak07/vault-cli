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
	return EncryptAndUpload(file, cfg, db)
}

func DownloadHandler(file string, cfg *Config, db *sql.DB) error {
	fmt.Printf("Downloading %s using mode=%s\n", file, cfg.Mode)
	if cfg.Mode == "local" {
		return LocalDownload(file)
	}
	return DownloadAndDecrypt(file, cfg, db)
}
