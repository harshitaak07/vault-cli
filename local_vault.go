package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func LocalUpload(filePath string) error {
	destDir := os.Getenv("VAULT_REMOTE_PATH")
	if destDir == "" {
		return fmt.Errorf("VAULT_REMOTE_PATH not set")
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(destDir, filepath.Base(filePath))
	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}

func LocalDownload(fileName string) error {
	srcDir := os.Getenv("VAULT_REMOTE_PATH")
	if srcDir == "" {
		return fmt.Errorf("VAULT_REMOTE_PATH not set")
	}
	src := filepath.Join(srcDir, filepath.Base(fileName))
	dest := "downloaded_" + filepath.Base(fileName)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}
