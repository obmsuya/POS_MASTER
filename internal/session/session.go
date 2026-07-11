package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Cached struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	FullName     string    `json:"full_name"`
	IsSuperuser  bool      `json:"is_superuser"`
	UserType     string    `json:"user_type"`
	SavedAt      time.Time `json:"saved_at"`
}

const maxAge = 12 * time.Hour

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".faltasi")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "session.json"), nil
}

func Load() (*Cached, error) {
	filePath, err := path()
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil
	}
	var cached Cached
	if err := json.Unmarshal(raw, &cached); err != nil {
		return nil, nil
	}
	if time.Since(cached.SavedAt) > maxAge {
		return nil, nil
	}
	return &cached, nil
}

func Save(cached Cached) error {
	filePath, err := path()
	if err != nil {
		return err
	}
	cached.SavedAt = time.Now()
	raw, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, raw, 0600)
}

func Clear() error {
	filePath, err := path()
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
