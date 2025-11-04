package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Session struct {
	User      string    `json:"user"`
	StartedAt time.Time `json:"started_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".vault", "session.json"), nil
}

func EnsureDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(home, ".vault"), 0700)
}

func Save(user string, ttl time.Duration) error {
	if err := EnsureDir(); err != nil {
		return err
	}
	p, _ := path()
	s := Session{
		User: user, StartedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(ttl),
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(p, b, 0600)
}

func Load() (*Session, error) {
	p, _ := path()
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if time.Now().UTC().After(s.ExpiresAt) {
		return nil, errors.New("session expired")
	}
	return &s, nil
}

func Clear() error {
	p, _ := path()
	_ = os.Remove(p)
	return nil
}

func Require() error {
	_, err := Load()
	return err
}
