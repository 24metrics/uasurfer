package uasurfer

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	amazonFireFingerprint = regexp.MustCompile(`\s(k[a-z]{3,5}|sd\d{4}ur)\s`) //tablet or phone
)

func (u *UserAgent) evalOS(ua string) bool {
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

	case specs == "macintosh":
		u.evalMacintosh(ua)

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
