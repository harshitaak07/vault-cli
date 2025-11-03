package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func EncryptAndUpload(filePath string, cfg *Config, db *sql.DB) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var plainKey, encryptedKey []byte
	if cfg.Mode == "local" {
		plainKey, encryptedKey, err = LocalKey()
		if err != nil {
			return err
		}
	} else {
		plainKey, encryptedKey, err = GenerateDataKey(cfg.KmsKey)
		if err != nil {
			return err
		}
	}

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, data, nil)

	encodedKey := base64.StdEncoding.EncodeToString(encryptedKey)
	hash := sha256Sum(data)
	keyName := filepath.Base(filePath)

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(awsCfg)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(keyName),
		Body:   bytes.NewReader(ciphertext),
		Metadata: map[string]string{
			"encryptedkey":    encodedKey,
			"encryption_mode": cfg.Mode,
			"uploader":        os.Getenv("VAULT_USER_ID"),
			"upload_ts":       time.Now().UTC().Format(time.RFC3339),
			"file_hash":       hash,
		},
	})
	if err != nil {
		return fmt.Errorf("s3 put: %w", err)
	}

	if db != nil {
		info, _ := os.Stat(filePath)
		RecordFile(db, keyName, hash, info.Size(), "s3", cfg.Mode)
		RecordAudit(db, "upload", keyName, "s3", true, "")
	}
	return nil
}

func DownloadAndDecrypt(fileName string, cfg *Config, db *sql.DB) error {
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(awsCfg)
	out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(filepath.Base(fileName)),
	})
	if err != nil {
		if db != nil {
			RecordAudit(db, "download", fileName, "s3", false, err.Error())
		}
		return err
	}
	defer out.Body.Close()
	data, err := io.ReadAll(out.Body)
	if err != nil {
		return err
	}
	meta := out.Metadata
	encodedKey := meta["encryptedkey"]
	if encodedKey == "" {
		return fmt.Errorf("missing encrypted key metadata")
	}
	encryptedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return err
	}

	var plainKey []byte
	if meta["encryption_mode"] == "local" || cfg.Mode == "local" {
		plainKey = encryptedKey
	} else {
		plainKey, err = DecryptDataKey(encryptedKey)
		if err != nil {
			return err
		}
	}

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	if fileHash, ok := meta["file_hash"]; ok {
		if sha256Sum(plaintext) != fileHash {
			return fmt.Errorf("hash mismatch: integrity check failed")
		}
	}

	outFile := "decrypted_" + filepath.Base(fileName)
	if err := os.WriteFile(outFile, plaintext, 0644); err != nil {
		return err
	}

	if db != nil {
		RecordAudit(db, "download", fileName, "s3", true, "")
	}
	return nil
}
