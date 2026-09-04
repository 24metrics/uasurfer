package uasurfer

import "strings"

// VersionFrozen marks an OS version token that the browser freezes across
// releases. Browser.Version.Extra keeps its existing version-suffix meaning.
const VersionFrozen = "frozen"

// IsFrozen reports whether Parse(ua, true) marked this OS version as frozen.
func (v Version) IsFrozen() bool { return v.Extra == VersionFrozen }

// OSVersionDetails describes the OS version stated by an agent without
// replacing it with an inferred value.
type OSVersionDetails struct {
	Reported Version
	Minimum  Version
	Frozen   bool
}

func firstFlag(flags []bool) bool { return len(flags) > 0 && flags[0] }

func resolveFrozenOS(rawUserAgent string, parsed *UserAgent) {
	ua := normalise(rawUserAgent)
	switch parsed.OS.Name {
	case OSiOS, OSiPadOS:
		resolveFrozenIOSVersion(ua, platformGroup(ua), &parsed.OS.Version)
	case OSMacOSX:
		markFrozenMacOSVersion(ua, &parsed.OS.Version)
	}
}

// resolveFrozenIOSVersion replaces a frozen iOS-family version with full
// Safari's release where possible and marks it otherwise.
func resolveFrozenIOSVersion(ua, agentPlatform string, version *Version) {
	if !hasAnyFrozenIOSVersionToken(agentPlatform) {
		return
	}

	candidate := UserAgent{}
	candidate.parseBrowserName(ua)
	candidate.parseBrowserVersion(ua)
	if candidate.Browser.Name == BrowserSafari && candidate.Browser.Version.Major >= 26 {
		candidate.Browser.Version.Extra = ""
		*version = candidate.Browser.Version
		return
	}

	version.Extra = VersionFrozen
}

func markFrozenMacOSVersion(ua string, version *Version) {
	candidate := UserAgent{}
	candidate.parseBrowserName(ua)
	candidate.parseBrowserVersion(ua)
	if hasFrozenMacOSVersion(ua, candidate.Browser) {
		version.Extra = VersionFrozen
	}
}

// IsFrozenOSVersion reports whether the agent's OS version is a known frozen
// browser token rather than a reliable measurement of the current OS.
func IsFrozenOSVersion(rawUserAgent string) bool {
	ua, parsed := parseOSVersionInput(rawUserAgent)
	return isFrozenOSVersion(ua, parsed)
}

func isFrozenOSVersion(ua string, parsed *UserAgent) bool {
	switch parsed.OS.Name {
	case OSiOS, OSiPadOS:
		return hasAnyFrozenIOSVersionToken(ua)
	case OSMacOSX:
		return hasFrozenMacOSVersion(ua, parsed.Browser)
	default:
		return false
	}
}

var safariMacOSFloors = []struct {
	safariMajor int
	macOS       Version
}{
	{16, Version{Major: 11}},
	{17, Version{Major: 12}},
	{18, Version{Major: 13}},
	{26, Version{Major: 14}},
}

// MinimumMacOSVersion returns the oldest macOS release supported by the
// Safari version in the agent, or the zero Version when no floor is known.
func MinimumMacOSVersion(rawUserAgent string) Version {
	ua, parsed := parseOSVersionInput(rawUserAgent)
	return minimumMacOSVersion(ua, parsed)
}

func minimumMacOSVersion(ua string, parsed *UserAgent) Version {
	if parsed.OS.Name != OSMacOSX || parsed.Browser.Name != BrowserSafari || strings.Contains(ua, "mobile/") {
		return Version{}
	}

	floor := Version{}
	for _, entry := range safariMacOSFloors {
		if parsed.Browser.Version.Major >= entry.safariMajor {
			floor = entry.macOS
		}
	}
	return floor
}

// ParseOSVersionDetails returns the stated version, its frozen status, and any
// conservative macOS floor without changing the reported value.
func ParseOSVersionDetails(rawUserAgent string) OSVersionDetails {
	ua, parsed := parseOSVersionInput(rawUserAgent)
	return OSVersionDetails{
		Reported: parsed.OS.Version,
		Minimum:  minimumMacOSVersion(ua, parsed),
		Frozen:   isFrozenOSVersion(ua, parsed),
	}
}

func parseOSVersionInput(rawUserAgent string) (string, *UserAgent) {
	ua := normalise(rawUserAgent)
	parsed := new(UserAgent)
	parse(ua, nil, parsed)
	return ua, parsed
}

func hasAnyFrozenIOSVersionToken(ua string) bool {
	return hasFrozenIOSVersionToken(ua, "18_6") ||
		hasFrozenIOSVersionToken(ua, "18_6_2") ||
		hasFrozenIOSVersionToken(ua, "18_7")
}

func hasFrozenMacOSVersion(ua string, browser Browser) bool {
	if macOSVersionToken(ua) != "10_15_7" {
		return false
	}

	switch browser.Name {
	case BrowserSafari:
		return browser.Version.Major >= 16
	case BrowserChrome:
		return browser.Version.Major >= 129
	default:
		return false
	}
}

func macOSVersionToken(ua string) string {
	start := strings.IndexByte(ua, '(')
	if start == -1 {
		return ""
	}
	end := strings.IndexByte(ua[start+1:], ')')
	if end == -1 {
		return ""
	}
	platform := ua[start+1 : start+1+end]
	const prefix = "mac os x "
	index := strings.Index(platform, prefix)
	if index == -1 {
		return ""
	}
	fields := strings.Fields(platform[index+len(prefix):])
	if len(fields) == 0 {
		return ""
	}
	return strings.TrimRight(fields[0], ";")
}

func hasFrozenIOSVersionToken(agentPlatform, version string) bool {
	return strings.Contains(agentPlatform, "cpu iphone os "+version+" like mac os x") ||
		strings.Contains(agentPlatform, "cpu os "+version+" like mac os x")
}
