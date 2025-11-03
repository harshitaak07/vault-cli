package main

import (
	"errors"
	"os"
)

type Config struct {
	Bucket          string
	KmsKey          string
	Region          string
	Mode            string // "kms" or "local"
	LocalPath       string
	RequirePassword bool
	PasswordFile    string
	DBPath          string
}

func LoadConfig() (*Config, error) {
	mode := os.Getenv("VAULT_MODE")
	if mode == "" {
		mode = "kms"
	}
	db := os.Getenv("VAULT_DB_PATH")
	if db == "" {
		db = "vault.db"
	}
	cfg := &Config{
		Bucket:       os.Getenv("VAULT_BUCKET"),
		KmsKey:       os.Getenv("VAULT_KMS_KEY"),
		Region:       os.Getenv("AWS_REGION"),
		Mode:         mode,
		LocalPath:    os.Getenv("VAULT_REMOTE_PATH"),
		PasswordFile: os.Getenv("VAULT_PASS_FILE"),
		DBPath:       db,
	}
	if os.Getenv("VAULT_REQUIRE_PASSWORD") == "1" {
		cfg.RequirePassword = true
	}
	if cfg.Mode == "kms" {
		if cfg.Bucket == "" || cfg.KmsKey == "" {
			return nil, errors.New("VAULT_BUCKET and VAULT_KMS_KEY must be set for kms mode")
		}
	}
	if cfg.Mode == "local" && cfg.LocalPath == "" {
		return nil, errors.New("VAULT_REMOTE_PATH must be set for local mode")
	}
	if cfg.RequirePassword && cfg.PasswordFile == "" {
		return nil, errors.New("VAULT_PASS_FILE must be set when VAULT_REQUIRE_PASSWORD=1")
	}
	return cfg, nil
}
