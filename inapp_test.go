package uasurfer

import "testing"

func TestInAppBrowsers(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		browser Browser
	}{
		{
			name:    "Line appended to full iOS Safari",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/604.1 Line/10.9.1",
			browser: Browser{BrowserLine, Version{10, 9, 1, ""}},
		},
		{
			name:    "Line Android webview",
			agent:   "Mozilla/5.0 (Linux; Android 13; SM-S901B Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.128 Mobile Safari/537.36 Line/13.1.0/IAB",
			browser: Browser{BrowserLine, Version{13, 1, 0, ""}},
		},
		{
			name:    "Facebook iOS",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBAV/456.0.0.36.107;FBBV/577025896]",
			browser: Browser{BrowserFacebook, Version{456, 0, 0, "36.107"}},
		},
		{
			name:    "Facebook Android",
			agent:   "Mozilla/5.0 (Linux; Android 5.1.1; R7sf Build/LMY47V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.149 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/263.0.0.46.121;]",
			browser: Browser{BrowserFacebook, Version{263, 0, 0, "46.121"}},
		},
		{
			name:    "Instagram iOS",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 289.0.0.13.62 (iPhone14,2; iOS 16_5; en_US)",
			browser: Browser{BrowserInstagram, Version{289, 0, 0, "13.62"}},
		},
		{
			name:    "WeChat Android",
			agent:   "Mozilla/5.0 (Linux; Android 11; SM-G991B Build/RP1A.200720.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.181 Mobile Safari/537.36 MicroMessenger/8.0.2.1860(0x28000234) Process/toolsmp WeChat/arm64",
			browser: Browser{BrowserWeChat, Version{8, 0, 2, "1860"}},
		},
		{
			name:    "TikTok iOS",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 musical_ly_30.5.0 JsSdk/2.0 NetType/WIFI",
			browser: Browser{BrowserTikTok, Version{30, 5, 0, ""}},
		},
		{
			name:    "TikTok Android Bytedance webview",
			agent:   "Mozilla/5.0 (Linux; Android 10; SM-A505F; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/95.0.4638.74 Mobile Safari/537.36 BytedanceWebview/d8a21c6",
			browser: Browser{BrowserTikTok, Version{}},
		},
		{
			name:    "Snapchat iOS",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Snapchat/12.10.0.36 (like Safari/8614.2.9.0.10)",
			browser: Browser{BrowserSnapchat, Version{12, 10, 0, "36"}},
		},
		{
			name:    "DuckDuckGo iOS",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 DuckDuckGo/7 Safari/605.1.15",
			browser: Browser{BrowserDuckDuckGo, Version{7, 0, 0, ""}},
		},
		{
			name:    "DuckDuckGo abbreviated marker",
			agent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15 Ddg/17.2",
			browser: Browser{BrowserDuckDuckGo, Version{17, 2, 0, ""}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Parse(test.agent).Browser; got != test.browser {
				t.Fatalf("Browser = %#v, want %#v", got, test.browser)
			}
		})
	}
}

func TestAppBrowserAnchoring(t *testing.T) {
	tests := []struct {
		name  string
		agent string
		want  BrowserName
	}{
		{
			name:  "baseline is not Line",
			agent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 baseline/2.0",
			want:  BrowserUnknown,
		},
		{
			name:  "glued Snapchat marker",
			agent: "Mozilla/5.0 (Linux; Android 12; SM-G991U; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.61 Mobile Safari/537.36snapchat11.79.0.33",
			want:  BrowserChrome,
		},
		{
			name:  "Instagram prefix is not an app field",
			agent: "Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Mobile Safari/537.36 InstagramHelper/1.0",
			want:  BrowserChrome,
		},
		{
			name:  "Facebook prefix is not an app field",
			agent: "Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Mobile Safari/537.36 FBAvailable/1.0",
			want:  BrowserChrome,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Parse(test.agent).Browser.Name; got != test.want {
				t.Fatalf("Browser.Name = %v, want %v", got, test.want)
			}
		})
	}

	agent := "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 baseline/2.0 Line/13.1.0"
	if got, want := Parse(agent).Browser, (Browser{BrowserLine, Version{13, 1, 0, ""}}); got != want {
		t.Fatalf("anchored app version = %#v, want %#v", got, want)
	}
}

// WebKit appends the client application name right after the engine token, so
// an app that replaces the default "Mobile/<build>" name instead of appending to
// it must still be identified.
func TestInAppBrowsersWithoutMobileToken(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		browser Browser
	}{
		{
			name:    "Instagram without Mobile token",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Instagram 300.0.0.13.62",
			browser: Browser{BrowserInstagram, Version{300, 0, 0, "13.62"}},
		},
		{
			name:    "WeChat without Mobile token",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) MicroMessenger/8.0.60",
			browser: Browser{BrowserWeChat, Version{8, 0, 60, ""}},
		},
		{
			name:    "bare WKWebView stays unidentified",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)",
			browser: Browser{BrowserUnknown, Version{}},
		},
		{
			name:    "baseline is not Line without Mobile token",
			agent:   "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) baseline/2.0",
			browser: Browser{BrowserUnknown, Version{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Parse(test.agent).Browser; got != test.browser {
				t.Fatalf("Browser = %#v, want %#v", got, test.browser)
			}
		})
	}
}

// An in-app browser that appends its name to full Safari's user agent must not be
// credited with Safari's release, because the app, not Safari, decides what it
// reports. With the heuristic enabled the stated numbers are kept and marked as
// frozen; Parse passes the token through unmarked.
func TestInAppFrozenIOSVersion(t *testing.T) {
	agent := "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.7 Mobile/15E148 Safari/604.1 Line/15.0.0"

	plain := Parse(agent)
	if plain.Browser.Name != BrowserLine {
		t.Fatalf("Browser.Name = %v, want %v", plain.Browser.Name, BrowserLine)
	}
	if want := (Version{18, 7, 0, ""}); plain.OS.Version != want {
		t.Fatalf("Parse OS.Version = %#v, want %#v", plain.OS.Version, want)
	}

	resolved := Parse(agent, true)
	if resolved.Browser.Name != BrowserLine {
		t.Fatalf("Browser.Name = %v, want %v", resolved.Browser.Name, BrowserLine)
	}
	if want := (Version{18, 7, 0, VersionFrozen}); resolved.OS.Version != want {
		t.Fatalf("Parse OS.Version = %#v, want %#v", resolved.OS.Version, want)
	}
	if !resolved.OS.Version.IsFrozen() {
		t.Fatal("OS.Version.IsFrozen() = false, want true")
	}
}

func TestInAppDetectionPreservesBots(t *testing.T) {
	tests := []struct {
		agent string
		want  BrowserName
	}{
		{"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1)", BrowserGoogleBot},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.2.5 (KHTML, like Gecko) Version/8.0.2 Safari/600.2.5 (Applebot/0.1)", BrowserAppleBot},
		{"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)", BrowserFacebookBot},
		{"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)", BrowserDuckDuckGoBot},
		{"Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0 Mobile Safari/537.36 Line/13.1.0 facebookexternalhit/1.1", BrowserFacebookBot},
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Snapchat/12.10.0 DuckDuckBot/1.0", BrowserDuckDuckGoBot},
	}

	for _, test := range tests {
		if got := Parse(test.agent).Browser.Name; got != test.want {
			t.Errorf("Browser.Name = %v, want %v for %q", got, test.want, test.agent)
		}
	}
}

func TestInAppBrowserStrings(t *testing.T) {
	tests := map[BrowserName]string{
		BrowserFacebook:   "BrowserFacebook",
		BrowserInstagram:  "BrowserInstagram",
		BrowserWeChat:     "BrowserWeChat",
		BrowserTikTok:     "BrowserTikTok",
		BrowserSnapchat:   "BrowserSnapchat",
		BrowserLine:       "BrowserLine",
		BrowserDuckDuckGo: "BrowserDuckDuckGo",
	}

	for browser, want := range tests {
		if got := browser.String(); got != want {
			t.Errorf("BrowserName.String() = %q, want %q", got, want)
		}
	}
}
