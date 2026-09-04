// Command uaparse prints what this package makes of a single user agent.
//
// Usage:
//
//	go run ./cmd/uaparse '<user agent>'
//	echo '<user agent>' | go run ./cmd/uaparse
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/24metrics/uasurfer"
)

func main() {
	agent, err := readUserAgent()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Resolving frozen versions keeps the stated numbers and either replaces or
	// marks them, which is the more informative mode for a diagnostic tool.
	result := uasurfer.Parse(agent, true)

	osVersion := formatVersion(result.OS.Version)
	if floor := uasurfer.MinimumMacOSVersion(agent); floor.Major > 0 {
		osVersion += fmt.Sprintf(", min. macOS %d", floor.Major)
	}

	fmt.Printf("Browser:         %s\n", result.Browser.Name.StringTrimPrefix())
	fmt.Printf("Browser version: %s\n", formatVersion(result.Browser.Version))
	fmt.Printf("OS:              %s\n", result.OS.Name.StringTrimPrefix())
	fmt.Printf("OS version:      %s\n", osVersion)
	fmt.Printf("Platform:        %s\n", result.OS.Platform.StringTrimPrefix())
	fmt.Printf("Device:          %s\n", result.DeviceType.StringTrimPrefix())
}

func readUserAgent() (string, error) {
	if len(os.Args) > 1 {
		return strings.Join(os.Args[1:], " "), nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	agent := strings.TrimSpace(string(data))
	if agent == "" {
		return "", fmt.Errorf("pass a user agent as an argument or on stdin")
	}
	return agent, nil
}

func formatVersion(version uasurfer.Version) string {
	if version == (uasurfer.Version{}) {
		return "unknown"
	}

	number := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
	switch {
	case version.IsFrozen():
		return number + "-" + uasurfer.VersionFrozen
	case version.Extra != "":
		return number + "." + version.Extra
	default:
		return number
	}
}
