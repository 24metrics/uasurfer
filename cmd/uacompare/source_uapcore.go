package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const uapCoreRawBase = "https://raw.githubusercontent.com/ua-parser/uap-core/master/tests/"

func fetchUAPCore(dir string) error {
	client := newClient()
	for _, name := range []string{"test_ua.yaml", "test_os.yaml"} {
		if err := download(client, uapCoreRawBase+name, filepath.Join(dir, name)); err != nil {
			return err
		}
	}
	return nil
}

// loadUAPCore reads the browser and OS fixtures. uap-core keeps one oracle per
// dimension in its own file, so a user agent appears in both with only the part
// that file is about.
func loadUAPCore(dir string) ([]expectation, error) {
	browsers, err := loadUAPCoreFile(filepath.Join(dir, "test_ua.yaml"))
	if err != nil {
		return nil, err
	}
	systems, err := loadUAPCoreFile(filepath.Join(dir, "test_os.yaml"))
	if err != nil {
		return nil, err
	}

	cases := make([]expectation, 0, len(browsers)+len(systems))
	for _, c := range browsers {
		exp := expectation{
			ua:             c.ua,
			browserLabel:   c.family,
			browserVersion: c.version(),
		}
		if names, ok := uapCoreBrowsers[c.family]; ok {
			exp.browsers = names
			exp.browserKnown = true
		}
		cases = append(cases, exp)
	}
	for _, c := range systems {
		exp := expectation{
			ua:        c.ua,
			osLabel:   c.family,
			osVersion: c.version(),
		}
		if names, ok := uapCoreSystems[c.family]; ok {
			exp.oses = names
			exp.osKnown = true
		}
		cases = append(cases, exp)
	}
	return cases, nil
}

// uapCoreCase is one uap-core fixture entry.
type uapCoreCase struct {
	ua     string
	family string
	major  string
	minor  string
	patch  string
}

func (c uapCoreCase) version() string {
	parts := []string{}
	for _, p := range []string{c.major, c.minor, c.patch} {
		if p == "" {
			break
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, ".")
}

// loadUAPCoreFile reads the uap-core fixture format, a flat list of single-level
// mappings.
func loadUAPCoreFile(path string) ([]uapCoreCase, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var (
		cases   []uapCoreCase
		current *uapCoreCase
	)
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "test_cases:" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			if current != nil {
				cases = append(cases, *current)
			}
			current = &uapCoreCase{}
			trimmed = strings.TrimSpace(trimmed[2:])
		}
		if current == nil {
			continue
		}

		key, value, ok := splitKeyValue(trimmed)
		if !ok {
			continue
		}
		switch key {
		case "user_agent_string":
			current.ua = value
		case "family":
			current.family = value
		case "major":
			current.major = value
		case "minor":
			current.minor = value
		case "patch":
			current.patch = value
		}
	}
	if current != nil {
		cases = append(cases, *current)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("%s: no test cases found", path)
	}
	return cases, nil
}
