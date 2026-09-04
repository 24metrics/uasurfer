package uasurfer

import "strings"

var appMarkers = []struct {
	marker string
	name   BrowserName
}{
	{"fbav/", BrowserFacebook},
	{"fban/", BrowserFacebook},
	{"fb_iab/", BrowserFacebook},
	{"instagram ", BrowserInstagram},
	{"micromessenger/", BrowserWeChat},
	{"bytedancewebview/", BrowserTikTok},
	{"musical_ly_", BrowserTikTok},
	{"trill_", BrowserTikTok},
	{"snapchat/", BrowserSnapchat},
	{"line/", BrowserLine},
	{"duckduckgo/", BrowserDuckDuckGo},
	{"ddg/", BrowserDuckDuckGo},
}

// appMarkerFirstBytes holds the first byte of every marker so a position that
// cannot start a marker is skipped without running a comparison per marker.
var appMarkerFirstBytes = func() (first [256]bool) {
	for _, app := range appMarkers {
		first[app.marker[0]] = true
	}
	return
}()

// appBrowser returns an app named in s. Markers only count at the start of a
// field so, for example, "baseline/" is not mistaken for Line.
func appBrowser(s string) BrowserName {
	for i := 0; i < len(s); i++ {
		if !appMarkerFirstBytes[s[i]] {
			continue
		}
		if i > 0 && !isAppEdge(s[i-1]) {
			continue
		}
		for _, app := range appMarkers {
			if strings.HasPrefix(s[i:], app.marker) {
				return app.name
			}
		}
	}
	return BrowserUnknown
}

func isAppEdge(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\v', '\f', ';', '(', '[', '/', ',', '-':
		return true
	default:
		return false
	}
}

func appVersion(v *Version, ua, prefix string) bool {
	tail := appIdentityTail(ua)
	index := anchoredAppMarker(tail, prefix)
	if index == -1 {
		return false
	}

	version := tail[index+len(prefix):]
	end := 0
	for end < len(version) {
		c := version[end]
		if (c < '0' || c > '9') && c != '.' {
			break
		}
		end++
	}
	version = strings.TrimRight(version[:end], ".")
	parts := strings.Split(version, ".")
	if len(parts) > 3 {
		v.Extra = strings.Join(parts[3:], ".")
	}
	return v.parse(version)
}

func appIdentityTail(ua string) string {
	if index := strings.Index(ua, "chrome/"); index != -1 {
		return ua[index+len("chrome/"):]
	}
	if tail := webkitAppTail(ua); tail != "" {
		return tail
	}
	// An app that sends only its own identity has no engine token to anchor on.
	return ua
}

func anchoredAppMarker(s, marker string) int {
	for offset := 0; offset < len(s); {
		index := strings.Index(s[offset:], marker)
		if index == -1 {
			return -1
		}
		index += offset
		if index == 0 || isAppEdge(s[index-1]) {
			return index
		}
		offset = index + 1
	}
	return -1
}

// webkitApp returns the app rendering the page on an Apple engine.
func webkitApp(ua string) BrowserName {
	return appBrowser(webkitAppTail(ua))
}

// webkitAppTail returns the part of an Apple WebKit user agent that carries the
// client's own identity, so unrelated product text cannot match a marker.
//
// WebKit builds its user agent by appending the client application name right
// after the engine token, and WKWebView passes "Mobile/<build>" as that name.
// Preferring the Mobile/ token keeps the search short for the common case, and
// falling back to the engine token still finds apps that replace the default
// name instead of appending to it.
func webkitAppTail(ua string) string {
	if index := strings.Index(ua, "mobile/"); index != -1 {
		return ua[index+len("mobile/"):]
	}
	if index := strings.Index(ua, "like gecko)"); index != -1 {
		return ua[index+len("like gecko)"):]
	}
	return ""
}

// chromiumBrowser searches only behind Chrome's own token. Branded browsers
// and Android in-app webviews append their identity in this part of the UA.
func chromiumBrowser(ua string) BrowserName {
	index := strings.Index(ua, "chrome/")
	if index == -1 {
		return BrowserChrome
	}
	if name := appBrowser(ua[index+len("chrome/"):]); name != BrowserUnknown {
		return name
	}
	return BrowserChrome
}
