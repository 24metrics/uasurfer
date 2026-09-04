// Package uasurfer provides fast and reliable abstraction
// of HTTP User-Agent strings. The philosophy is to identify
// technologies that holds >1% market share, and to avoid
// expending resources and accuracy on guessing at esoteric UA
// strings.
package uasurfer

import "strings"

//go:generate stringer -type=DeviceType,BrowserName,OSName,Platform -output=const_string.go

// DeviceType (int) returns a constant.
type DeviceType int

// A complete list of supported devices in the
// form of constants.
const (
	DeviceUnknown DeviceType = iota
	DeviceComputer
	DeviceTablet
	DevicePhone
	DeviceConsole
	DeviceWearable
	DeviceTV
	DeviceMediaHub
)

// StringTrimPrefix is like String() but trims the "Device" prefix
func (d DeviceType) StringTrimPrefix() string {
	return strings.TrimPrefix(d.String(), "Device")
}

// BrowserName (int) returns a constant.
type BrowserName int

// A complete list of supported web browsers in the
// form of constants.
const (
	BrowserUnknown BrowserName = iota
	BrowserChrome
	BrowserIE
	BrowserSafari
	BrowserFirefox
	BrowserAndroid
	BrowserOpera
	BrowserBlackberry
	BrowserUCBrowser
	BrowserSilk
	BrowserNokia
	BrowserNetFront
	BrowserQQ
	BrowserMaxthon
	BrowserSogouExplorer
	BrowserSpotify
	BrowserNintendo
	BrowserSamsung
	BrowserYandex
	BrowserCocCoc
	BrowserBot // Bot list begins here
	BrowserAppleBot
	BrowserBaiduBot
	BrowserBingBot
	BrowserDuckDuckGoBot
	BrowserFacebookBot
	BrowserGoogleBot
	BrowserLinkedInBot
	BrowserMsnBot
	BrowserPingdomBot
	BrowserTwitterBot
	BrowserYandexBot
	BrowserCocCocBot
	BrowserYahooBot // Bot list ends here
	BrowserFacebook
	BrowserInstagram
	BrowserWeChat
	BrowserTikTok
	BrowserSnapchat
	BrowserLine
	BrowserDuckDuckGo
)

// StringTrimPrefix is like String() but trims the "Browser" prefix
func (b BrowserName) StringTrimPrefix() string {
	return strings.TrimPrefix(b.String(), "Browser")
}

// OSName (int) returns a constant.
type OSName int

// A complete list of supported OSes in the
// form of constants. For handling particular versions
// of operating systems (e.g. Windows 2000), see
// the README.md file.
const (
	OSUnknown OSName = iota
	OSWindowsPhone
	OSWindows
	OSMacOSX
	OSiOS
	OSAndroid
	OSBlackberry
	OSChromeOS
	OSKindle
	OSWebOS
	OSLinux
	OSPlaystation
	OSXbox
	OSNintendo
	OSBot
)

// StringTrimPrefix is like String() but trims the "OS" prefix
func (o OSName) StringTrimPrefix() string {
	return strings.TrimPrefix(o.String(), "OS")
}

// Platform (int) returns a constant.
type Platform int

// A complete list of supported platforms in the
// form of constants. Many OSes report their
// true platform, such as Android OS being Linux
// platform.
const (
	PlatformUnknown Platform = iota
	PlatformWindows
	PlatformMac
	PlatformLinux
	PlatformiPad
	PlatformiPhone
	PlatformiPod
	PlatformBlackberry
	PlatformWindowsPhone
	PlatformPlaystation
	PlatformXbox
	PlatformNintendo
	PlatformBot
)

// StringTrimPrefix is like String() but trims the "Platform" prefix
func (p Platform) StringTrimPrefix() string {
	return strings.TrimPrefix(p.String(), "Platform")
}

type Version struct {
	Major int
	Minor int
	Patch int
	Extra string
}

// VersionFrozen is the marker Parse(ua, true) assigns to OS.Version.Extra when
// it keeps an OS version that the browser reports as a constant for every
// release. The numbers are the ones stated in the user agent, so they remain
// usable, while the marker says they cannot be read as the version of the system
// in front of the user.
//
// Browser.Version.Extra retains its existing meaning: components beyond the
// patch level. IsFrozen should therefore be used on an OS version produced by
// Parse(ua, true), not as a general classifier for arbitrary Version values.
// A macOS floor is not encoded in Extra; MinimumMacOSVersion returns it as a
// separate Version without changing the parsed UserAgent.
const VersionFrozen = "frozen"

// IsFrozen reports whether this version is a constant the browser reports for
// every release rather than the version of the system, which is the case for the
// iOS-family and desktop platform tokens that Apple, Google and Mozilla cap.
//
// Only Parse(ua, true) sets the marker. Without the flag Parse states what the
// user agent says, so use IsFrozenOSVersion there.
func (v Version) IsFrozen() bool {
	return v.Extra == VersionFrozen
}

func (v Version) Less(c Version) bool {
	if v.Major < c.Major {
		return true
	}

	if v.Major > c.Major {
		return false
	}

	if v.Minor < c.Minor {
		return true
	}

	if v.Minor > c.Minor {
		return false
	}

	return v.Patch < c.Patch
}

type UserAgent struct {
	Browser    Browser
	OS         OS
	DeviceType DeviceType
}

type Browser struct {
	Name    BrowserName
	Version Version
}

type OS struct {
	Platform Platform
	Name     OSName
	Version  Version
}

// OSVersionDetails describes the OS-version information that can be derived
// from a user agent without replacing the value it actually reports.
type OSVersionDetails struct {
	// Reported is the OS version stated in the user agent, identical to
	// Parse(rawUA).OS.Version. It never carries VersionFrozen.
	Reported Version

	// Minimum is a conservative macOS floor derived from Safari's version, or
	// the zero Version when no floor can be established.
	Minimum Version

	// Frozen reports whether Reported is a known constant rather than a
	// reliable measurement of the current OS version.
	Frozen bool
}

// Reset resets the UserAgent to it's zero value
func (ua *UserAgent) Reset() {
	ua.Browser = Browser{}
	ua.OS = OS{}
	ua.DeviceType = DeviceUnknown
}

// IsBot returns true if the UserAgent represent a bot
func (ua *UserAgent) IsBot() bool {
	if ua.Browser.Name >= BrowserBot && ua.Browser.Name <= BrowserYahooBot {
		return true
	}
	if ua.OS.Name == OSBot {
		return true
	}
	if ua.OS.Platform == PlatformBot {
		return true
	}
	return false
}

// Parse accepts a raw user agent (string) and returns the UserAgent.
//
// Called with one argument it reports what the user agent states and applies no
// heuristics, so an OS version that the browser reports as a frozen constant is
// passed through unchanged. Use IsFrozenOSVersion to recognise such a value.
//
// The optional resolveFrozenOSVersion asks for that constant to be dealt with
// instead of passed through. An iOS user agent whose platform token is frozen then
// takes its OS version from full Safari's own version, which tracks the iOS
// release, and keeps the stated numbers with the VersionFrozen marker where no such
// hint exists, as in a WKWebView, an in-app browser or a third-party browser. On
// macOS nothing can be recovered, so the value is only marked. The result is then
// no longer a plain reading of the user agent, which is why this has to be asked
// for.
//
//	Parse(ua)        // as stated
//	Parse(ua, false) // identical
//	Parse(ua, true)  // resolved where possible, marked otherwise
//
// Only the first value is read; any further ones are ignored.
func Parse(ua string, resolveFrozenOSVersion ...bool) *UserAgent {
	dest := new(UserAgent)
	parse(ua, dest, firstFlag(resolveFrozenOSVersion))
	return dest
}

// ParseUserAgent is the same as Parse, but populates the supplied UserAgent.
// It is the caller's responsibility to call Reset() on the UserAgent before
// passing it to this function.
func ParseUserAgent(ua string, dest *UserAgent, resolveFrozenOSVersion ...bool) {
	parse(ua, dest, firstFlag(resolveFrozenOSVersion))
}

func firstFlag(flags []bool) bool {
	return len(flags) > 0 && flags[0]
}

func parse(ua string, dest *UserAgent, resolveFrozenOSVersion bool) {
	ua = normalise(ua)
	switch {
	case len(ua) == 0:
		dest.OS.Platform = PlatformUnknown
		dest.OS.Name = OSUnknown
		dest.Browser.Name = BrowserUnknown
		dest.DeviceType = DeviceUnknown

	// stop on on first case returning true
	case dest.evalOS(ua, resolveFrozenOSVersion):
	case dest.evalBrowserName(ua):
	default:
		dest.evalBrowserVersion(ua)
		dest.evalDevice(ua)
	}
}

// normalise normalises the user supplied agent string so that
// we can more easily parse it.
func normalise(ua string) string {
	if len(ua) <= 1024 {
		var buf [1024]byte
		ascii := copyLower(buf[:len(ua)], ua)
		if !ascii {
			// Fall back for non ascii characters
			return strings.ToLower(ua)
		}
		return string(buf[:len(ua)])
	}
	// Fallback for unusually long strings
	return strings.ToLower(ua)
}

// copyLower copies a lowercase version of s to b. It assumes s contains only single byte characters
// and will panic if b is nil or is not long enough to contain all the bytes from s.
// It returns early with false if any characters were non ascii.
func copyLower(b []byte, s string) bool {
	for j := 0; j < len(s); j++ {
		c := s[j]
		if c > 127 {
			return false
		}

		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}

		b[j] = c
	}
	return true
}
