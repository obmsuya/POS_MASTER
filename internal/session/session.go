// Package session caches a logged-in operator's JWT pair locally so
// relaunching the CLI mid-shift doesn't require typing a password every
// time. Nothing about the password itself is ever stored.
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

// maxAge bounds how long a cached session is trusted before requiring a
// fresh login, independent of whether the refresh token itself is still
// technically valid server-side — keeps a lost/shared laptop from staying
// logged in indefinitely.
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

// Load returns the cached session if one exists and hasn't expired.
func Load() (*Cached, error) {
	filePath, err := path()
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil // no session file yet — not an error
	}
	var cached Cached
	if err := json.Unmarshal(raw, &cached); err != nil {
		return nil, nil // corrupt file — treat as no session
	}
	if time.Since(cached.SavedAt) > maxAge {
		return nil, nil
	}
	return &cached, nil
}

// Save persists a fresh session to disk with the current timestamp.
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

// Clear removes any cached session (used on logout or a rejected token).
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
