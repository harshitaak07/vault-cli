package core

import (
	"database/sql"
	"fmt"
	"time"

	"vault-cli/internal/aws"
	"vault-cli/internal/config"
	"vault-cli/internal/tui"
)

func UploadHandler(file string, cfg *config.Config, db *sql.DB) error {
	tui.ShowVaultBanner()

	fmt.Println("\nStarting Secure Upload Process...\n")
	fmt.Println("Features:")
	fmt.Println("AES-256 encryption with AWS KMS key")
	fmt.Println("Secure S3 upload over TLS")
	fmt.Println("Local SQLite metadata tracking")
	fmt.Println("CloudWatch audit logging\n")

	start := time.Now()
	go tui.RunTUI()
	err := aws.EncryptAndUpload(file, cfg, db)
	if err != nil {
		fmt.Printf("\n❌ Upload Failed: %v\n", err)
		return fmt.Errorf("encrypt/upload: %w", err)
	}
	_ = aws.LogToCloudWatch(fmt.Sprintf("File uploaded: %s", file))

	elapsed := time.Since(start)

	fmt.Printf("\nUpload Complete!\n")
	fmt.Printf("File: %s\n", file)
	fmt.Printf("Hash: SHA256_PLACEHOLDER\n")
	fmt.Printf("Storage: AWS S3\n")
	fmt.Printf("Mode: %s\n", cfg.Mode)
	fmt.Printf("Duration: %v\n", elapsed)

	return nil
}

func DownloadHandler(file string, cfg *config.Config, db *sql.DB) error {
	tui.ShowVaultBanner()

	fmt.Println("Features:")
	fmt.Println("  • Fetch encrypted file from AWS S3")
	fmt.Println("  • Decrypt via AWS KMS data key")
	fmt.Println("  • Integrity verification via SHA-256")

	start := time.Now()
	go tui.RunTUI()
	err := aws.DownloadAndDecrypt(file, cfg, db)
	if err != nil {
		fmt.Printf("\n❌ Download Failed: %v\n", err)
		return fmt.Errorf("decrypt/download: %w", err)
	}
	_ = aws.LogToCloudWatch(fmt.Sprintf("File downloaded: %s", file))

	elapsed := time.Since(start)

	fmt.Printf("\nDownload Complete!\n")
	fmt.Printf("File: %s\n", file)
	fmt.Printf("Hash: SHA256_PLACEHOLDER\n")
	fmt.Printf("Saved As: decrypted_%s\n", file)
	fmt.Printf("Mode: %s\n", cfg.Mode)
	fmt.Printf("Duration: %v\n", elapsed)

	return nil
}
