package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/24metrics/uasurfer"
)

// expectation is one reference case, already translated into this package's
// vocabulary. The label fields keep the original reference names so a difference
// can be declared per family rather than per user agent.
type expectation struct {
	ua string

	browserLabel   string
	browsers       []uasurfer.BrowserName
	browserKnown   bool
	browserVersion string

	osLabel   string
	oses      []uasurfer.OSName
	osKnown   bool
	osVersion string

	deviceLabel string
	devices     []uasurfer.DeviceType
	deviceKnown bool
}

func compare(cases []expectation, declared deviationSet) *report {
	browserName := &dimension{name: "Browser name"}
	browserVersion := &dimension{name: "Browser version"}
	osName := &dimension{name: "OS name"}
	osVersion := &dimension{name: "OS version"}
	deviceType := &dimension{name: "Device type"}

	for _, c := range cases {
		got := uasurfer.Parse(c.ua)

		// Browser.
		switch {
		case !c.browserKnown:
			browserName.skipped++
			browserVersion.skipped++
		case declared.covers("browser", c.ua, c.browserLabel):
			browserName.known++
			browserVersion.known++
		case containsBrowser(c.browsers, got.Browser.Name):
			browserName.hit()
			compareVersion(browserVersion, declared, "browser_version", c.ua, c.browserLabel, c.browserVersion, got.Browser.Version)
		default:
			browserName.miss(c.ua, joinBrowsers(c.browsers), got.Browser.Name.String())
			// A version is only meaningful once the product itself matches.
			browserVersion.skipped++
		}

		// Operating system.
		switch {
		case !c.osKnown:
			osName.skipped++
			osVersion.skipped++
		case declared.covers("os", c.ua, c.osLabel):
			osName.known++
			osVersion.known++
		case containsOS(c.oses, got.OS.Name):
			osName.hit()
			compareVersion(osVersion, declared, "os_version", c.ua, c.osLabel, c.osVersion, got.OS.Version)
		default:
			osName.miss(c.ua, joinOSes(c.oses), got.OS.Name.String())
			osVersion.skipped++
		}

		// Device type.
		switch {
		case !c.deviceKnown:
			deviceType.skipped++
		case declared.covers("device", c.ua, c.deviceLabel):
			deviceType.known++
		case containsDevice(c.devices, got.DeviceType):
			deviceType.hit()
		default:
			deviceType.miss(c.ua, joinDevices(c.devices), got.DeviceType.String())
		}
	}

	dims := []*dimension{browserName, browserVersion, osName, osVersion}
	if deviceType.compared > 0 || deviceType.known > 0 {
		dims = append(dims, deviceType)
	}
	return &report{dimensions: dims}
}

func compareVersion(d *dimension, declared deviationSet, dimensionName, ua, label, want string, got uasurfer.Version) {
	stated := statedParts(want)
	if len(stated) == 0 {
		d.skipped++
		return
	}
	if declared.covers(dimensionName, ua, label) {
		d.known++
		return
	}
	if versionMatches(stated, got) {
		d.hit()
		return
	}
	d.miss(ua, want, fmt.Sprintf("%d.%d.%d", got.Major, got.Minor, got.Patch))
}

// statedParts returns the leading numeric parts a reference actually states.
// References omit trailing parts a user agent does not carry, and some use
// labels such as a build name, which cannot be compared against integers.
func statedParts(version string) []int {
	version = strings.TrimSpace(version)
	if version == "" {
		return nil
	}
	var parts []int
	for i, field := range strings.Split(version, ".") {
		if i == 3 {
			break
		}
		n, err := strconv.Atoi(field)
		if err != nil {
			break
		}
		parts = append(parts, n)
	}
	return parts
}

func versionMatches(stated []int, got uasurfer.Version) bool {
	actual := []int{got.Major, got.Minor, got.Patch}
	for i, want := range stated {
		if actual[i] != want {
			return false
		}
	}
	return true
}

func containsBrowser(allowed []uasurfer.BrowserName, got uasurfer.BrowserName) bool {
	for _, a := range allowed {
		if a == got {
			return true
		}
	}
	return false
}

func containsOS(allowed []uasurfer.OSName, got uasurfer.OSName) bool {
	for _, a := range allowed {
		if a == got {
			return true
		}
	}
	return false
}

func containsDevice(allowed []uasurfer.DeviceType, got uasurfer.DeviceType) bool {
	for _, a := range allowed {
		if a == got {
			return true
		}
	}
	return false
}

func joinBrowsers(names []uasurfer.BrowserName) string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = n.StringTrimPrefix()
	}
	return strings.Join(out, "|")
}

func joinOSes(names []uasurfer.OSName) string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = n.StringTrimPrefix()
	}
	return strings.Join(out, "|")
}

func joinDevices(names []uasurfer.DeviceType) string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = n.StringTrimPrefix()
	}
	return strings.Join(out, "|")
}

// deviationSet holds the differences that are intentional. A deviation is
// declared either for a whole reference label or for a single user agent.
type deviationSet struct {
	labels     map[string]map[string]bool
	agents     map[string]map[string]bool
	substrings map[string]map[string]bool
}

func (s deviationSet) covers(dimension, ua, label string) bool {
	if label != "" {
		if dims, ok := s.labels[label]; ok && dims[dimension] {
			return true
		}
	}
	if dims, ok := s.agents[ua]; ok && dims[dimension] {
		return true
	}
	for substring, dims := range s.substrings {
		if dims[dimension] && strings.Contains(strings.ToLower(ua), substring) {
			return true
		}
	}
	return false
}

func loadDeviations(path string) (deviationSet, error) {
	set := deviationSet{
		labels:     map[string]map[string]bool{},
		agents:     map[string]map[string]bool{},
		substrings: map[string]map[string]bool{},
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return set, nil
	}
	if err != nil {
		return set, err
	}

	for i, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			return set, fmt.Errorf("%s:%d: want scope, dimension, key and reason separated by tabs", path, i+1)
		}
		scope, dimension, key := fields[0], fields[1], fields[2]

		var target map[string]map[string]bool
		switch scope {
		case "label":
			target = set.labels
		case "ua":
			target = set.agents
		case "contains":
			target = set.substrings
			key = strings.ToLower(key)
		default:
			return set, fmt.Errorf("%s:%d: unknown scope %q", path, i+1, scope)
		}
		if target[key] == nil {
			target[key] = map[string]bool{}
		}
		target[key][dimension] = true
	}
	return set, nil
}
