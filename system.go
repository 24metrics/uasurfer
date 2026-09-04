package uasurfer

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	amazonFireFingerprint = regexp.MustCompile(`\s(k[a-z]{3,5}|sd\d{4}ur)\s`) //tablet or phone
)

func (u *UserAgent) evalOS(ua string, resolveFrozenOSVersion bool) bool {
	// This is a subjective parsing for cfnetwork user agents.
	// For us, we consider this as iOS from in-app browsers.
	if strings.Contains(ua, "cfnetwork/") {
		u.OS.Platform = PlatformUnknown
		u.OS.Name = OSiOS
		return u.maybeBot()
	}

	s := strings.IndexRune(ua, '(')
	e := strings.IndexRune(ua, ')')
	if s > e {
		s = 0
		e = len(ua)
	}
	if e == -1 {
		e = len(ua)
	}

	agentPlatform := ua[s+1 : e]
	specsEnd := strings.Index(agentPlatform, ";")
	var specs string
	if specsEnd != -1 {
		specs = agentPlatform[:specsEnd]
	} else {
		specs = agentPlatform
	}

	//strict OS & version identification
	switch {
	case specs == "android":
		u.evalLinux(ua, agentPlatform)

	case specs == "bb10" || specs == "playbook":
		u.OS.Platform = PlatformBlackberry
		u.OS.Name = OSBlackberry

	case specs == "x11" || specs == "linux":
		u.evalLinux(ua, agentPlatform)

	case strings.HasPrefix(specs, "ipad") || strings.HasPrefix(specs, "iphone") || strings.HasPrefix(specs, "ipod touch") || strings.HasPrefix(specs, "ipod"):
		u.evaliOS(specs, agentPlatform)
		if resolveFrozenOSVersion {
			resolveFrozenIOSVersion(ua, agentPlatform, &u.OS.Version)
		}

	case specs == "macintosh":
		u.evalMacintosh(ua)
		if resolveFrozenOSVersion {
			markFrozenMacOSVersion(ua, &u.OS.Version)
		}

	default:
		switch {
		// Blackberry
		case strings.Contains(ua, "blackberry") || strings.Contains(ua, "playbook"):
			u.OS.Platform = PlatformBlackberry
			u.OS.Name = OSBlackberry

		// Windows Phone
		case strings.Contains(agentPlatform, "windows phone ") &&
			!strings.Contains(agentPlatform, "xbox"): // Xbox one user agents have 'windows phone'
			u.evalWindowsPhone(agentPlatform)

		// Windows, Xbox
		case strings.Contains(ua, "windows ") || strings.Contains(ua, "microsoft-cryptoapi"):
			u.evalWindows(ua)

		// Kindle
		case strings.Contains(ua, "kindle/") || amazonFireFingerprint.MatchString(agentPlatform):
			u.OS.Platform = PlatformLinux
			u.OS.Name = OSKindle

		// Linux (broader attempt)
		case strings.Contains(ua, "linux") || strings.Contains(ua, "crkey"):
			u.evalLinux(ua, agentPlatform)

		// WebOS (non-linux flagged)
		case strings.Contains(ua, "webos") || strings.Contains(ua, "hpwos"):
			u.OS.Platform = PlatformLinux
			u.OS.Name = OSWebOS

		// Nintendo
		case strings.Contains(ua, "nintendo"):
			u.OS.Platform = PlatformNintendo
			u.OS.Name = OSNintendo

		// Playstation
		case strings.Contains(ua, "playstation") || strings.Contains(ua, "vita") || strings.Contains(ua, "psp"):
			u.OS.Platform = PlatformPlaystation
			u.OS.Name = OSPlaystation

		// Android
		case strings.Contains(ua, "android"):
			u.evalLinux(ua, agentPlatform)

		// Apple CFNetwork
		case strings.Contains(ua, "cfnetwork") && strings.Contains(ua, "darwin"):
			u.evalMacintosh(ua)

		default:
			u.OS.Platform = PlatformUnknown
			u.OS.Name = OSUnknown
		}
	}

	return u.maybeBot()
}

// maybeBot checks if the UserAgent is a bot and sets
// all bot related fields if it is
func (u *UserAgent) maybeBot() bool {
	if u.IsBot() {
		u.OS.Platform = PlatformBot
		u.OS.Name = OSBot
		u.DeviceType = DeviceComputer
		return true
	}
	return false
}

// evalLinux returns the `Platform`, `OSName` and Version of UAs with
// 'linux' listed as their platform.
func (u *UserAgent) evalLinux(ua string, agentPlatform string) {

	switch {
	// Kindle Fire
	case strings.Contains(ua, "kindle") || amazonFireFingerprint.MatchString(agentPlatform):
		// get the version of Android if available, though we don't call this OSAndroid
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSKindle
		u.OS.Version.findVersionNumber(agentPlatform, "android ")

	// Android, Kindle Fire
	case strings.Contains(ua, "android") ||
		strings.Contains(ua, "googletv") ||
		strings.Contains(ua, "crkey"):
		// Android
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSAndroid
		u.OS.Version.findVersionNumber(agentPlatform, "android ")

	// ChromeOS
	case strings.Contains(ua, "cros"):
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSChromeOS

	// WebOS
	case strings.Contains(ua, "webos") || strings.Contains(ua, "hpwos"):
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSWebOS

	// Linux, "Linux-like"
	case strings.Contains(ua, "x11") || strings.Contains(ua, "bsd") || strings.Contains(ua, "suse") || strings.Contains(ua, "debian") || strings.Contains(ua, "ubuntu"):
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSLinux

	default:
		u.OS.Platform = PlatformLinux
		u.OS.Name = OSLinux
	}
}

// evaliOS returns the `Platform`, `OSName` and Version of UAs with
// 'ipad' or 'iphone' listed as their platform.
func (u *UserAgent) evaliOS(uaPlatform string, agentPlatform string) {

	switch {
	// iPhone
	case strings.HasPrefix(uaPlatform, "iphone"):
		u.OS.Platform = PlatformiPhone
		u.OS.Name = OSiOS
		u.OS.getiOSVersion(agentPlatform)

	// iPad
	case strings.HasPrefix(uaPlatform, "ipad"):
		u.OS.Platform = PlatformiPad
		u.OS.Name = OSiOS
		u.OS.getiOSVersion(agentPlatform)

	// iPod
	case strings.HasPrefix(uaPlatform, "ipod touch") || strings.HasPrefix(uaPlatform, "ipod"):
		u.OS.Platform = PlatformiPod
		u.OS.Name = OSiOS
		u.OS.getiOSVersion(agentPlatform)

	default:
		u.OS.Platform = PlatformiPad
		u.OS.Name = OSUnknown
	}
}

func (u *UserAgent) evalWindowsPhone(agentPlatform string) {
	u.OS.Platform = PlatformWindowsPhone

	if u.OS.Version.findVersionNumber(agentPlatform, "windows phone os ") || u.OS.Version.findVersionNumber(agentPlatform, "windows phone ") {
		u.OS.Name = OSWindowsPhone
	} else {
		u.OS.Name = OSUnknown
	}
}

func (u *UserAgent) evalWindows(ua string) {

	switch {
	//Xbox -- it reads just like Windows
	case strings.Contains(ua, "xbox") || strings.Contains(ua, "Xbox One"):
		u.OS.Platform = PlatformXbox
		u.OS.Name = OSXbox
		if !u.OS.Version.findVersionNumber(ua, "windows phone ") {
			if !u.OS.Version.findVersionNumber(ua, "windows nt ") {
				u.OS.Version.Major = 6
				u.OS.Version.Minor = 0
				u.OS.Version.Patch = 0
			}
		}

	// No windows version
	case !strings.Contains(ua, "windows "):
		u.OS.Platform = PlatformWindows
		u.OS.Name = OSUnknown

	case strings.Contains(ua, "windows nt ") && u.OS.Version.findVersionNumber(ua, "windows nt "):
		u.OS.Platform = PlatformWindows
		u.OS.Name = OSWindows

	case strings.Contains(ua, "windows xp"):
		u.OS.Platform = PlatformWindows
		u.OS.Name = OSWindows
		u.OS.Version.Major = 5
		u.OS.Version.Minor = 1
		u.OS.Version.Patch = 0

	default:
		u.OS.Platform = PlatformWindows
		u.OS.Name = OSUnknown

	}
}

func (u *UserAgent) evalMacintosh(uaPlatformGroup string) {
	u.OS.Platform = PlatformMac
	if i := strings.Index(uaPlatformGroup, "os x "); i != -1 {
		u.OS.Name = OSMacOSX
		u.OS.Version.parse(uaPlatformGroup[i+5:])

		return
	}
	u.OS.Name = OSUnknown
}

func (v *Version) findVersionNumber(s string, versionPrefix string) bool {
	if ind := strings.Index(s, versionPrefix); ind != -1 {
		// Given `s` as "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101.5 Safari/537.36",
		// and `versionPrefix` as "Chrome/",
		// then `ind` is 75,
		// then `s[ind+len(versionPrefix):]` will become "30.0.1599.101.5 Safari/537.36"
		versionString := s[ind+len(versionPrefix):]
		// Trim any words after the version number.
		// For example, given "30.0.1599.101.5 Safari/537.36", we want to trim it to "30.0.1599.101.5"
		{
			ss := strings.Split(versionString, " ")
			if len(ss) > 0 {
				versionString = ss[0]
			}
		}
		// Trim the version string of any trailing characters that are not part of the version number
		versionString = strings.TrimRight(versionString, ":;.,/\\")
		{
			ss := strings.Split(versionString, ".")
			// If there are more than 3 parts to the version number, then the last part is the label.
			if len(ss) > 3 {
				v.Extra = strings.Join(ss[3:], ".")
			}
		}
		return v.parse(versionString)
	}
	return false
}

// getiOSVersion accepts the platform portion of a UA string and returns
// a Version.
func (o *OS) getiOSVersion(uaPlatformGroup string) {
	if i := strings.Index(uaPlatformGroup, "cpu iphone os "); i != -1 {
		o.Version.parse(uaPlatformGroup[i+14:])
		return
	}

	if i := strings.Index(uaPlatformGroup, "cpu os "); i != -1 {
		o.Version.parse(uaPlatformGroup[i+7:])
		return
	}

	o.Version.parse(uaPlatformGroup)
}

// resolveFrozenIOSVersion replaces a frozen iOS-family version with the release
// it can be recovered from, and marks it when it cannot. The numbers stated in the
// user agent are kept either way, so a caller never loses the value; the marker
// says it is a constant rather than a measurement.
func resolveFrozenIOSVersion(ua, agentPlatform string, version *Version) {
	if !hasAnyFrozenIOSVersionToken(agentPlatform) {
		return
	}

	// Safari is tied to iOS releases, so its Version/ token is the best
	// available iOS version once the platform token is frozen. Reuse the normal
	// browser classifier so WKWebViews, in-app browsers and named third-party
	// browsers cannot be mistaken for full Safari.
	candidate := UserAgent{}
	candidate.evalBrowserName(ua)
	if candidate.Browser.Name == BrowserSafari {
		var safariVersion Version
		if safariVersion.findVersionNumber(ua, "version/") && safariVersion.Major >= 26 {
			// Extra belongs to Safari's version token. The OS version has only
			// three numeric components; VersionFrozen is its sole Extra value.
			safariVersion.Extra = ""
			*version = safariVersion
			return
		}
	}

	version.Extra = VersionFrozen
}

// markFrozenMacOSVersion flags a capped desktop version. Nothing can be recovered
// here, because one Safari release serves several macOS versions and the other
// browsers cap the token without offering a usable version of their own, so the
// value is only marked. See MinimumMacOSVersion for the lower bound that can still
// be derived.
func markFrozenMacOSVersion(ua string, version *Version) {
	candidate := UserAgent{}
	candidate.evalBrowserName(ua)
	candidate.evalBrowserVersion(ua)

	if hasFrozenMacOSVersion(ua, candidate.Browser) {
		version.Extra = VersionFrozen
	}
}

// IsFrozenOSVersion reports whether the OS version Parse returns for this user
// agent is a constant the browser sends for every release, rather than the
// system in front of the user. When it is true, treat OS.Version as unknown.
//
// Apple and Google cap the platform token they report. A Mac running macOS 26
// and one running macOS 11 therefore send the same value. The browser version
// proves the value cannot be genuine only after that browser drops Catalina:
//
//	Safari  16 and later, since Catalina supports no Safari beyond 15.6
//	Chrome  129 and later, since Catalina supports no Chrome beyond 128
//
// Firefox has capped its token at 10.15 since version 87, but still supports
// Catalina. Its token may therefore be genuine and is deliberately not marked.
//
// On the iOS family WebKit hardcodes the token for every client, so any of the
// frozen values counts, whatever browser sent it. Note that those values are also
// real iOS releases: a device genuinely running 18.7 is reported as frozen too,
// because the two cannot be told apart.
//
// The value is about the user agent, not about a particular call: for an iOS user
// agent that states full Safari 26 or later, Parse(ua, true) can replace the frozen
// token with the release recovered from the Safari version, and this function
// still reports the token itself as frozen.
//
// A false result means the version is either genuine or not provably frozen, not
// that it is guaranteed to be accurate.
func IsFrozenOSVersion(rawUserAgent string) bool {
	ua, parsed := parseOSVersionInput(rawUserAgent)
	return isFrozenOSVersion(ua, parsed)
}

func isFrozenOSVersion(ua string, parsed *UserAgent) bool {
	switch parsed.OS.Name {
	case OSiOS:
		return hasAnyFrozenIOSVersionToken(ua)

	case OSMacOSX:
		return hasFrozenMacOSVersion(ua, parsed.Browser)

	default:
		return false
	}
}

// safariMacOSFloors maps a Safari major release to the oldest macOS version it
// was offered for. Safari on the Mac ships for the current system and the two
// before it, so its version does not identify a single release, but it does rule
// out everything older than the entry here.
//
// The values come from the overview of Apple's Safari release notes, for example
// Safari 26 being available for macOS 26, macOS Sequoia and macOS Sonoma.
//
// Releases before 16 are absent on purpose: their floor is at or below the
// capped token itself, so it would add nothing. There is no formula behind these
// numbers, and the jump from macOS 15 to macOS 26 shows why, so a new Safari
// generation needs a new entry here.
var safariMacOSFloors = []struct {
	safariMajor int
	macOS       Version
}{
	{16, Version{Major: 11}}, // Big Sur, Monterey, Ventura
	{17, Version{Major: 12}}, // Monterey, Ventura, Sonoma
	{18, Version{Major: 13}}, // Ventura, Sonoma, macOS 15
	{26, Version{Major: 14}}, // Sonoma, Sequoia, macOS 26
}

// MinimumMacOSVersion returns the oldest macOS release that could be behind this
// user agent, or the zero Version when no bound can be derived. The result is
// separate metadata: it is not stored in OS.Version.Extra and this function does
// not modify a UserAgent returned by Parse.
//
// It exists because the macOS version in a user agent is capped and therefore
// useless on its own, while Safari's own version still carries information: a
// given Safari release runs on at most three macOS versions, so the oldest of
// them is a floor that cannot be wrong. Safari 26.5 means macOS 14 or newer.
//
// The result is deliberately a floor, not a guess. Deriving an exact version by
// subtracting a constant from the Safari version is wrong twice over: it ignores
// that Safari ships for two older systems, and the year-based renumbering in 2025
// broke any fixed offset.
//
// The zero Version is returned for anything that carries no such hint: other
// browsers, which cap the token without a usable version of their own, and Safari
// before 16, whose floor would be older than the capped token anyway. A Safari
// generation newer than the newest known entry falls back to that entry, which
// stays a valid floor.
//
// An iPad requesting a desktop site can send a UA indistinguishable from a
// Mac. The result is therefore a macOS floor only if the request actually came
// from a Mac. If a Mobile/ token makes the iPad origin visible, zero is returned.
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

// ParseOSVersionDetails returns the OS version stated in rawUserAgent together
// with any independently derivable macOS floor and its frozen-token status.
//
// Reported is identical to Parse(rawUserAgent).OS.Version: it is never resolved
// from Safari's version and never receives VersionFrozen in Extra. Use
// Parse(rawUserAgent, true) when the opt-in resolved-or-marked representation is
// wanted instead. Minimum and Frozen have the same semantics as
// MinimumMacOSVersion and IsFrozenOSVersion.
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
	parse(ua, parsed, false)
	return ua, parsed
}

func hasAnyFrozenIOSVersionToken(ua string) bool {
	return hasFrozenIOSVersionToken(ua, "18_6") ||
		hasFrozenIOSVersionToken(ua, "18_6_2") ||
		hasFrozenIOSVersionToken(ua, "18_7")
}

func hasFrozenMacOSVersion(ua string, browser Browser) bool {
	// WebKit and Chromium use exactly this underscore-separated cap. Read it
	// from the platform group rather than matching arbitrary text elsewhere.
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

// strToInt simply accepts a string and returns a `int`,
// with '0' being default.
func strToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

// parse accepts a string and returns a Version,
// with {0, 0, 0, ""} being default.
func (v *Version) parse(str string) bool {
	// This checks if the string is empty or if the first character of the string is not a digit.
	if len(str) == 0 || str[0] < '0' || str[0] > '9' {
		return false
	}

	// This is a for loop that iterates three times, with i starting from 0 and incrementing by 1 each time.
	// Within this loop, another for loop is started that iterates over each character c in the str string, along with its index k.
	for i := 0; i < 3; i++ {
		empty := true
		val := 0
		strLen := len(str) - 1

		for k, c := range str {
			// If the current character c is a digit, then the if empty block is executed.
			// If empty is true, val is set to the numeric value of c and empty is set to false.
			if c >= '0' && c <= '9' {
				if empty {
					val = int(c) - 48
					empty = false
					// If the current character is the last character in the string, then the str variable is set to an empty string.
					if k == strLen {
						str = str[:0]
					}
					continue
				}
				// then str is set to the substring starting from index k,
				// and the loop is exited using the break statement.

				// If val is zero and the current character is 0.
				if val == 0 {
					if c == '0' {
						// If the current character is the last character in the string, then str is set to an empty string.
						if k == strLen {
							str = str[:0]
						}
						continue
					}
					str = str[k:]
					break
				}

				// If val is not zero and the current character is a digit,
				// then val is updated by multiplying it by 10 and adding the numeric value of the current character.
				val = 10*val + int(c) - 48
				// If the current character is the last character in the string,
				// then str is set to an empty string.
				if k == strLen {
					str = str[:0]
				}
				continue
			}

			// From here, the current character is not a digit.
			// Else, the str is set to the substring starting from index k+1,
			str = str[k+1:]
			break
		}

		switch i {
		case 0:
			v.Major = val

		case 1:
			v.Minor = val

		case 2:
			v.Patch = val
		}
	}

	return true
}
