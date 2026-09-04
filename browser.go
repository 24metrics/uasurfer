package uasurfer

import (
	"strings"
)

// Browser struct contains the lowercase name of the browser, along
// with its browser version number. Browser are grouped together without
// consideration for device. For example, Chrome (Chrome/43.0) and Chrome for iOS
// (CriOS/43.0) would both return as "chrome" (name) and 43.0 (version). Similarly
// Internet Explorer 11 and Edge 12 would return as "ie" and "11" or "12", respectively.
// type Browser struct {
// 		Name    BrowserName
// 		Version struct {
// 				Major int
// 			Minor int
// 			Patch int
// 		}
// }

// Retrieve browser name from UA strings
func (u *UserAgent) evalBrowserName(ua string) bool {
	if bot := botBrowser(ua); bot != BrowserUnknown {
		u.Browser.Name = bot
		return u.maybeBot()
	}

	// Blackberry goes first because it reads as MSIE & Safari
	if strings.Contains(ua, "blackberry") || strings.Contains(ua, "playbook") || strings.Contains(ua, "bb10") || strings.Contains(ua, "rim ") {
		u.Browser.Name = BrowserBlackberry
		return u.maybeBot()
	}

	// https://user-agents.net/applications/dalvik
	if strings.Contains(ua, "dalvik/") {
		u.Browser.Name = BrowserAndroid
		return u.maybeBot()
	}

	// https://user-agents.net/string/cfnetwork-1126-darwin-19-5-0
	if strings.Contains(ua, "cfnetwork/") {
		u.Browser.Name = BrowserUnknown
		return u.maybeBot()
	}

	if strings.Contains(ua, "applewebkit") {
		inApp := webkitApp(ua)
		switch {
		case strings.Contains(ua, "googlebot"):
			u.Browser.Name = BrowserGoogleBot

		case strings.Contains(ua, "applebot"):
			u.Browser.Name = BrowserAppleBot

		case inApp != BrowserUnknown:
			u.Browser.Name = inApp

		case strings.Contains(ua, "qq/") || strings.Contains(ua, "qqbrowser/"):
			u.Browser.Name = BrowserQQ

		case strings.Contains(ua, "opr/") || strings.Contains(ua, "opios/"):
			u.Browser.Name = BrowserOpera

		case strings.Contains(ua, "silk/"):
			u.Browser.Name = BrowserSilk

		case strings.Contains(ua, "edg/") || strings.Contains(ua, "edgios/") || strings.Contains(ua, "edga/") || strings.Contains(ua, "edge/") || strings.Contains(ua, "iemobile/") || strings.Contains(ua, "msie "):
			u.Browser.Name = BrowserIE

		case strings.Contains(ua, "ucbrowser/") || strings.Contains(ua, "ucweb/"):
			u.Browser.Name = BrowserUCBrowser

		case strings.Contains(ua, "nintendobrowser/"):
			u.Browser.Name = BrowserNintendo

		case strings.Contains(ua, "samsungbrowser/"):
			u.Browser.Name = BrowserSamsung

		case strings.Contains(ua, "coc_coc_browser/"):
			u.Browser.Name = BrowserCocCoc

		// Yandex ships its engine under several names: the browser itself, the
		// browser inside its search app, and the search app's own webview.
		case strings.Contains(ua, "yabrowser/") || strings.Contains(ua, "yasearchbrowser/") ||
			strings.Contains(ua, "yandexsearchbrowser/") || strings.Contains(ua, "yandexsearch/") ||
			strings.Contains(ua, "yasearchapp/") || strings.Contains(ua, "yaapp_android/"):
			u.Browser.Name = BrowserYandex

		// Edge, Silk and other chrome-identifying browsers must evaluate before chrome, unless we want to add more overhead
		case strings.Contains(ua, "chrome/") || strings.Contains(ua, "crios/") || strings.Contains(ua, "chromium/") || strings.Contains(ua, "crmo/"):
			u.Browser.Name = chromiumBrowser(ua)

		case strings.Contains(ua, "android") && !strings.Contains(ua, "chrome/") && strings.Contains(ua, "version/") && !strings.Contains(ua, "like android"):
			// Android WebView on Android >= 4.4 is purposefully being identified as Chrome above -- https://developer.chrome.com/multidevice/webview/overview
			u.Browser.Name = BrowserAndroid

		case strings.Contains(ua, "fxios"):
			u.Browser.Name = BrowserFirefox

		case strings.Contains(ua, " spotify/"):
			u.Browser.Name = BrowserSpotify

		// Full Safari on Apple platforms includes both Version/ and Safari/.
		case strings.Contains(ua, "like gecko") && strings.Contains(ua, "mozilla/") && strings.Contains(ua, "version/") && strings.Contains(ua, "safari/") && !strings.Contains(ua, "linux") && !strings.Contains(ua, "android") && !strings.Contains(ua, "browser/") && !strings.Contains(ua, "os/") && !strings.Contains(ua, "yabrowser/"):
			u.Browser.Name = BrowserSafari

		default:
			goto notwebkit

		}
		return u.maybeBot()
	}

notwebkit:
	switch {
	case strings.Contains(ua, "nintendo 3ds") &&
		!strings.Contains(ua, "nintendobrowser/"):
		u.Browser.Name = BrowserNetFront

	case strings.Contains(ua, "qq/") || strings.Contains(ua, "qqbrowser/"):
		u.Browser.Name = BrowserQQ

	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		u.Browser.Name = BrowserIE

	case strings.Contains(ua, "gecko") && (strings.Contains(ua, "firefox") || strings.Contains(ua, "iceweasel") || strings.Contains(ua, "seamonkey") || strings.Contains(ua, "icecat")):
		u.Browser.Name = BrowserFirefox

	case strings.Contains(ua, "presto") || strings.Contains(ua, "opera"):
		u.Browser.Name = BrowserOpera

	case strings.Contains(ua, "ucbrowser"):
		u.Browser.Name = BrowserUCBrowser

	case strings.Contains(ua, "applebot"):
		u.Browser.Name = BrowserAppleBot

	case strings.Contains(ua, "baiduspider"):
		u.Browser.Name = BrowserBaiduBot

	case strings.Contains(ua, "adidxbot") || strings.Contains(ua, "bingbot") || strings.Contains(ua, "bingpreview"):
		u.Browser.Name = BrowserBingBot

	case strings.Contains(ua, "duckduckbot"):
		u.Browser.Name = BrowserDuckDuckGoBot

	case strings.Contains(ua, "facebot") || strings.Contains(ua, "facebookexternalhit"):
		u.Browser.Name = BrowserFacebookBot

	case strings.Contains(ua, "googlebot"):
		u.Browser.Name = BrowserGoogleBot

	case strings.Contains(ua, "linkedinbot"):
		u.Browser.Name = BrowserLinkedInBot

	case strings.Contains(ua, "msnbot"):
		u.Browser.Name = BrowserMsnBot

	case strings.Contains(ua, "pingdom.com_bot"):
		u.Browser.Name = BrowserPingdomBot

	case strings.Contains(ua, "twitterbot"):
		u.Browser.Name = BrowserTwitterBot

	case strings.Contains(ua, "yandex") || strings.Contains(ua, "yadirectfetcher"):
		u.Browser.Name = BrowserYandexBot

	case strings.Contains(ua, "yahoo"):
		u.Browser.Name = BrowserYahooBot

	case strings.Contains(ua, "coccocbot"):
		u.Browser.Name = BrowserCocCocBot

	case strings.Contains(ua, "phantomjs"):
		u.Browser.Name = BrowserBot

	// Some apps send only their own identity, with no engine token to anchor on.
	// This runs last so a named browser always wins.
	case appBrowser(ua) != BrowserUnknown:
		u.Browser.Name = appBrowser(ua)

	default:
		u.Browser.Name = BrowserUnknown

	}

	return u.maybeBot()
}

func botBrowser(ua string) BrowserName {
	switch {
	case strings.Contains(ua, "applebot"):
		return BrowserAppleBot
	case strings.Contains(ua, "baiduspider"):
		return BrowserBaiduBot
	case strings.Contains(ua, "adidxbot") || strings.Contains(ua, "bingbot") || strings.Contains(ua, "bingpreview"):
		return BrowserBingBot
	case strings.Contains(ua, "duckduckbot"):
		return BrowserDuckDuckGoBot
	case strings.Contains(ua, "facebot") || strings.Contains(ua, "facebookexternalhit"):
		return BrowserFacebookBot
	case strings.Contains(ua, "googlebot"):
		return BrowserGoogleBot
	case strings.Contains(ua, "linkedinbot"):
		return BrowserLinkedInBot
	case strings.Contains(ua, "msnbot"):
		return BrowserMsnBot
	case strings.Contains(ua, "pingdom.com_bot"):
		return BrowserPingdomBot
	case strings.Contains(ua, "twitterbot"):
		return BrowserTwitterBot
	// Only unambiguous crawler tokens belong here. A bare "yandex" or "yahoo"
	// also appears in the browsers and apps those vendors ship, so those stay in
	// the fallback below where a browser has already had its chance to match.
	case strings.Contains(ua, "yandexbot") || strings.Contains(ua, "yadirectfetcher"):
		return BrowserYandexBot
	case strings.Contains(ua, "yahoo! slurp") || strings.Contains(ua, "yahooseeker"):
		return BrowserYahooBot
	case strings.Contains(ua, "coccocbot"):
		return BrowserCocCocBot
	case strings.Contains(ua, "phantomjs"):
		return BrowserBot
	default:
		return BrowserUnknown
	}
}

// browserVersionTokens lists the version tokens a browser writes for itself,
// in the order they should be trusted. These take precedence over the generic
// "version/" token, because a browser running on an Apple or Android webview
// shares its user agent with the host: iOS adds Safari's "Version/17.5" and an
// Android webview adds "Version/4.0", neither of which is the browser's own
// version.
//
// Browsers whose real version only ever appears in "version/" are absent here:
// Safari, the stock Android browser, BlackBerry and Presto-era Opera.
var browserVersionTokens = map[BrowserName][]string{
	BrowserChrome:    {"crios/", "chrome/", "crmo/", "chromium/"},
	BrowserIE:        {"edgios/", "edga/", "edg/", "edge/", "msie "},
	BrowserFirefox:   {"fxios/", "firefox/"},
	BrowserOpera:     {"opios/", "opr/"},
	BrowserUCBrowser: {"ucbrowser/", "ucweb/"},
	BrowserQQ:        {"qq/", "qqbrowser/"},
	BrowserSamsung:   {"samsungbrowser/"},
	BrowserYandex: {
		"yabrowser/", "yasearchbrowser/", "yandexsearchbrowser/",
		"yandexsearch/", "yasearchapp/", "yaapp_android/",
	},
	BrowserCocCoc:   {"coc_coc_browser/"},
	BrowserSilk:     {"silk/"},
	BrowserSpotify:  {"spotify/"},
	BrowserNintendo: {"nintendobrowser/"},
}

// Retrieve browser version
// Methods used in order:
// 1st: use the browser's own version token (e.g. crios/123, ucbrowser/13)
// 2nd: look for generic version/#
// 3rd: derive it (MSIE from trident, Safari from the OS)
func (u *UserAgent) evalBrowserVersion(ua string) {
	// In-app browsers state their version in their own field only.
	switch u.Browser.Name {
	case BrowserFacebook:
		appVersion(&u.Browser.Version, ua, "fbav/")
		return
	case BrowserInstagram:
		appVersion(&u.Browser.Version, ua, "instagram ")
		return
	case BrowserWeChat:
		appVersion(&u.Browser.Version, ua, "micromessenger/")
		return
	case BrowserTikTok:
		if !appVersion(&u.Browser.Version, ua, "musical_ly_") {
			appVersion(&u.Browser.Version, ua, "trill_")
		}
		return
	case BrowserSnapchat:
		appVersion(&u.Browser.Version, ua, "snapchat/")
		return
	case BrowserLine:
		appVersion(&u.Browser.Version, ua, "line/")
		return
	case BrowserDuckDuckGo:
		if !appVersion(&u.Browser.Version, ua, "duckduckgo/") {
			appVersion(&u.Browser.Version, ua, "ddg/")
		}
		return
	}

	for _, token := range browserVersionTokens[u.Browser.Name] {
		if u.Browser.Version.findVersionNumber(ua, token) {
			return
		}
	}

	// if there is a 'version/#' attribute with numeric version, use it -- except for Chrome since Android vendors sometimes hijack version/#
	if u.Browser.Name != BrowserChrome && u.Browser.Version.findVersionNumber(ua, "version/") {
		return
	}

	switch u.Browser.Name {
	case BrowserAndroid:
		_ = u.Browser.Version.findVersionNumber(ua, "dalvik/")

	case BrowserOpera:
		// Presto-era Opera writes its release in opera/ and, on some builds, in
		// version/ handled above.
		_ = u.Browser.Version.findVersionNumber(ua, "opera/")

	case BrowserIE:
		// get MSIE version from trident version https://en.wikipedia.org/wiki/Trident_(layout_engine)
		if u.Browser.Version.findVersionNumber(ua, "trident/") {
			// convert trident versions 3-7 to MSIE version
			if (u.Browser.Version.Major >= 3) && (u.Browser.Version.Major <= 7) {
				u.Browser.Version.Major += 4
			}
		}

	case BrowserSafari: // executes typically if we're on iOS and not using a familiar browser
		u.Browser.Version = u.OS.Version
		// A marker belongs to the OS version, not to the browser version derived
		// from it.
		if u.Browser.Version.IsFrozen() {
			u.Browser.Version.Extra = ""
		}
		// early Safari used a version number +1 to OS version
		if (u.Browser.Version.Major <= 3) && (u.Browser.Version.Major >= 1) {
			u.Browser.Version.Major++
		}
	}
}
