package core

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"vault-cli/internal/aws"
	"vault-cli/internal/config"
	"vault-cli/internal/local"
	"vault-cli/internal/ui"
)

// UploadHandler — now with fancy progress UI
func UploadHandler(file string, cfg *config.Config, db *sql.DB) error {
	ui.ShowFeatureTable([]string{
		"Encrypt file",
		"Upload to S3",
		"Record to DB",
		"Audit to DynamoDB",
	})

	bar := ui.StartProgress("Uploading " + file)

	for i := 0; i < 50; i++ {
		bar.Add(1)
		time.Sleep(30 * time.Millisecond)
	}

	start := time.Now()
	if cfg.Mode == "local" {
		if err := local.LocalUpload(file); err != nil {
			ui.ShowErrorBox("❌ Upload Failed", err.Error())
			return fmt.Errorf("local upload: %w", err)
		}
		ui.ShowSuccessBox("✅ Upload Complete", fmt.Sprintf("%s uploaded locally in %v", file, time.Since(start)))
		info, _ := os.Stat(file)
		ui.ShowFileDetails(file, "N/A", info.Size(), "local", "filesystem")
		return nil
	}

	err := aws.EncryptAndUpload(file, cfg, db)
	if err != nil {
		ui.ShowErrorBox("❌ Upload Failed", err.Error())
		return fmt.Errorf("encrypt/upload: %w", err)
	}

	_ = aws.LogToCloudWatch(fmt.Sprintf("File uploaded: %s", file))

	ui.ShowSuccessBox("✅ Upload Complete", fmt.Sprintf("%s uploaded successfully in %v", file, time.Since(start)))

	info, _ := os.Stat(file)
	ui.ShowFileDetails(file, "SHA256_PLACEHOLDER", info.Size(), cfg.Mode, "s3")

	return nil
}

// DownloadHandler — now with fancy progress UI
func DownloadHandler(file string, cfg *config.Config, db *sql.DB) error {
	ui.ShowFeatureTable([]string{
		"Fetch file from S3",
		"Decrypt with KMS/local key",
		"Verify integrity",
		"Save decrypted copy",
	})

	bar := ui.StartProgress("Downloading " + file)
	for i := 0; i < 50; i++ {
		bar.Add(1)
		time.Sleep(30 * time.Millisecond)
	}

	start := time.Now()
	if cfg.Mode == "local" {
		if err := local.LocalDownload(file); err != nil {
			ui.ShowErrorBox("❌ Download Failed", err.Error())
			return fmt.Errorf("local download: %w", err)
		}
		ui.ShowSuccessBox("✅ Download Complete", fmt.Sprintf("%s downloaded locally in %v", file, time.Since(start)))
		ui.ShowFileDetails(file, "N/A", 0, "local", "filesystem")
		return nil
	}

	err := aws.DownloadAndDecrypt(file, cfg, db)
	if err != nil {
		ui.ShowErrorBox("❌ Download Failed", err.Error())
		return fmt.Errorf("decrypt/download: %w", err)
	}

	_ = aws.LogToCloudWatch(fmt.Sprintf("File downloaded: %s", file))
	ui.ShowSuccessBox("✅ Download Complete", fmt.Sprintf("%s downloaded successfully in %v", file, time.Since(start)))

	info, _ := os.Stat("decrypted_" + file)
	ui.ShowFileDetails(file, "SHA256_PLACEHOLDER", info.Size(), cfg.Mode, "s3")

	return nil
}
