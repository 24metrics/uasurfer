package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	matomoContentsURL = "https://api.github.com/repos/matomo-org/device-detector/contents/Tests/fixtures"
	matomoRawBase     = "https://raw.githubusercontent.com/matomo-org/device-detector/master/Tests/fixtures/"
)

// skippedMatomoFixtures are files whose schema is not the user agent shape this
// command reads: the client hint cases carry headers instead of a user agent,
// and the bot cases describe a crawler rather than a browser and an OS.
var skippedMatomoFixtures = map[string]bool{
	"clienthints.yml":     true,
	"clienthints-app.yml": true,
	"bots.yml":            true,
}

func fetchMatomo(dir string) error {
	client := newClient()
	resp, err := client.Get(matomoContentsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: unexpected status %s", matomoContentsURL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var entries []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &entries); err != nil {
		return err
	}

	fetched := 0
	for _, entry := range entries {
		if entry.Type != "file" || !strings.HasSuffix(entry.Name, ".yml") || skippedMatomoFixtures[entry.Name] {
			continue
		}
		if err := download(client, matomoRawBase+entry.Name, filepath.Join(dir, entry.Name)); err != nil {
			return err
		}
		fetched++
	}
	if fetched == 0 {
		return fmt.Errorf("no fixtures found at %s", matomoContentsURL)
	}
	fmt.Printf("Downloaded %d fixture files\n", fetched)
	return nil
}

func loadMatomo(dir string) ([]expectation, error) {
	names, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no *.yml fixtures in %s", dir)
	}

	var cases []expectation
	for _, name := range names {
		if skippedMatomoFixtures[filepath.Base(name)] {
			continue
		}
		data, err := os.ReadFile(name)
		if err != nil {
			return nil, err
		}
		cases = append(cases, parseMatomo(string(data))...)
	}
	return cases, nil
}

// matomoCase mirrors the fields this command reads from a fixture entry.
type matomoCase struct {
	ua            string
	osName        string
	osVersion     string
	clientName    string
	clientVersion string
	deviceType    string
}

// parseMatomo reads the Matomo fixture shape: a list of entries whose keys sit
// at one indentation level and whose os, client and device details sit one level
// deeper. The format is regular enough for a line reader, which keeps this
// command free of a YAML dependency.
func parseMatomo(data string) []expectation {
	var (
		cases   []expectation
		current *matomoCase
		section string
	)

	flush := func() {
		if current != nil && current.ua != "" {
			cases = append(cases, current.toExpectation())
		}
		current = nil
	}

	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(line, "---") {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 && strings.HasPrefix(trimmed, "-") {
			flush()
			current = &matomoCase{}
			section = ""
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			if trimmed == "" {
				continue
			}
			indent = 2
		}
		if current == nil {
			continue
		}

		key, value, ok := splitKeyValue(trimmed)
		if !ok {
			continue
		}
		if indent <= 2 {
			if value == "" {
				section = key
				continue
			}
			section = ""
			if key == "user_agent" {
				current.ua = value
			}
			continue
		}

		switch section + "." + key {
		case "os.name":
			current.osName = value
		case "os.version":
			current.osVersion = value
		case "client.name":
			current.clientName = value
		case "client.version":
			current.clientVersion = value
		case "device.type":
			current.deviceType = value
		}
	}
	flush()
	return cases
}

func (c matomoCase) toExpectation() expectation {
	exp := expectation{
		ua:             c.ua,
		browserLabel:   c.clientName,
		browserVersion: c.clientVersion,
		osLabel:        c.osName,
		osVersion:      c.osVersion,
		deviceLabel:    c.deviceType,
	}
	if names, ok := matomoClients[c.clientName]; ok {
		exp.browsers = names
		exp.browserKnown = true
	}
	if names, ok := matomoSystems[c.osName]; ok {
		exp.oses = names
		exp.osKnown = true
	}
	if types, ok := matomoDevices[c.deviceType]; ok {
		exp.devices = types
		exp.deviceKnown = true
	}
	return exp
}

// splitKeyValue splits a mapping line at its first colon. Values are unquoted,
// single-quoted or double-quoted, and user agents frequently contain colons of
// their own, so only the first one may be treated as the separator.
func splitKeyValue(line string) (key, value string, ok bool) {
	colon := strings.Index(line, ":")
	if colon == -1 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:colon])
	value = strings.TrimSpace(line[colon+1:])
	if len(value) >= 2 {
		switch value[0] {
		case '"':
			if value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
			}
		case '\'':
			if value[len(value)-1] == '\'' {
				value = strings.ReplaceAll(value[1:len(value)-1], "''", "'")
			}
		}
	} else if value == `""` || value == "''" {
		value = ""
	}
	return key, value, true
}
