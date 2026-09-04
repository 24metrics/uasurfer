package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// source is a reference corpus: where to get it and how to read it.
type source struct {
	name    string
	repo    string
	license string
	fetch   func(dir string) error
	load    func(dir string) ([]expectation, error)
}

func lookupSource(name string) (source, error) {
	switch name {
	case "matomo":
		return source{
			name:    "matomo",
			repo:    "github.com/matomo-org/device-detector",
			license: "LGPL-3.0-or-later, fetched locally, not redistributed here",
			fetch:   fetchMatomo,
			load:    loadMatomo,
		}, nil
	case "uap-core":
		return source{
			name:    "uap-core",
			repo:    "github.com/ua-parser/uap-core",
			license: "Apache-2.0",
			fetch:   fetchUAPCore,
			load:    loadUAPCore,
		}, nil
	default:
		return source{}, fmt.Errorf("unknown source %q, want matomo or uap-core", name)
	}
}

func download(client *http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: unexpected status %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, body, 0o644)
}

func newClient() *http.Client {
	return &http.Client{Timeout: 120 * time.Second}
}
