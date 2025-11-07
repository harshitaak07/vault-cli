package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"vault-cli/internal/aws"
	"vault-cli/internal/config"
	"vault-cli/internal/core"
)

type Secret struct {
	Category string
	Name     string
	Value    string
    CreatedAt string
    UpdatedAt string
}

func Add(database *sql.DB, cfg *config.Config, s Secret) error {
	plain := []byte(s.Value)
	var plainKey, wrappedKey []byte
	var err error

	if cfg.Mode == "local" {
		plainKey, wrappedKey, err = aws.LocalKey()
	} else {
		plainKey, wrappedKey, err = aws.GenerateDataKey(cfg.KmsKey)
	}
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}
	defer zero(plainKey)

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nil, nonce, plain, nil)

	hash := core.BytesSHA256Hex(plain)
	encodedWrapped := base64.StdEncoding.EncodeToString(wrappedKey)

	_, err = database.Exec(`
		INSERT INTO secrets(category, name, ciphertext, nonce, mode, hash, created_at, updated_at)
		VALUES(?,?,?,?,?,?,?,?) 
		ON CONFLICT(category, name) DO UPDATE SET 
			ciphertext=excluded.ciphertext,
			nonce=excluded.nonce,
			mode=excluded.mode,
			hash=excluded.hash,
			updated_at=excluded.updated_at
	`, s.Category, s.Name, encodedWrapped+"."+base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonce), cfg.Mode, hash, now(), now())
	return err
}

func Get(database *sql.DB, cfg *config.Config, category, name string) (string, error) {
	row := database.QueryRow(`SELECT ciphertext, nonce, mode FROM secrets WHERE category=? AND name=?`, category, name)
	var storedCT, nonceB64, mode string
	if err := row.Scan(&storedCT, &nonceB64, &mode); err != nil {
		return "", err
	}

	dot := -1
	for i := 0; i < len(storedCT); i++ {
		if storedCT[i] == '.' {
			dot = i
			break
		}
	}
	if dot < 0 {
		return "", fmt.Errorf("malformed ciphertext record")
	}
	wrappedB64 := storedCT[:dot]
	ctB64 := storedCT[dot+1:]

	wrappedKey, err := base64.StdEncoding.DecodeString(wrappedB64)
	if err != nil {
		return "", err
	}
	ct, err := base64.StdEncoding.DecodeString(ctB64)
	if err != nil {
		return "", err
	}
	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return "", err
	}

	var plainKey []byte
	if mode == "local" || cfg.Mode == "local" {
		plainKey = wrappedKey
	} else {
		plainKey, err = aws.DecryptDataKey(wrappedKey)
		if err != nil {
			return "", fmt.Errorf("decrypt data key: %w", err)
		}
	}
	defer zero(plainKey)

	block, err := aes.NewCipher(plainKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func List(database *sql.DB, category string) ([]Secret, error) {
    q := `SELECT category, name, created_at, updated_at FROM secrets`
	args := []any{}
	if category != "" {
		q += ` WHERE category=?`
		args = append(args, category)
	}
	q += ` ORDER BY category, name`
	rows, err := database.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Secret
    for rows.Next() {
        var c, n, created, updated string
        if err := rows.Scan(&c, &n, &created, &updated); err != nil {
			return nil, err
		}
        out = append(out, Secret{Category: c, Name: n, CreatedAt: created, UpdatedAt: updated})
	}
	return out, rows.Err()
}

func Delete(database *sql.DB, category, name string) error {
	_, err := database.Exec(`DELETE FROM secrets WHERE category=? AND name=?`, category, name)
	return err
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
func now() string { return time.Now().UTC().Format(time.RFC3339) }
