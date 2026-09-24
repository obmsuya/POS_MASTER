package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/inconshreveable/go-update"
)

const releaseAPI = "https://api.github.com/repos/obmsuya/POS_MASTER/releases/latest"

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func LatestRelease() (*Release, error) {
	resp, err := http.Get(releaseAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned %d", resp.StatusCode)
	}
	var r Release
	return &r, json.NewDecoder(resp.Body).Decode(&r)
}

func AssetURL(r *Release) string {
	name := fmt.Sprintf("faltasi-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	for _, a := range r.Assets {
		if a.Name == name {
			return a.BrowserDownloadURL
		}
	}
	return ""
}

func Apply(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return update.Apply(resp.Body, update.Options{})
}
