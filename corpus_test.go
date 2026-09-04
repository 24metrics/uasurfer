package uasurfer

import "testing"

// The cases below were found by comparing this package against an external
// corpus with cmd/uacompare. Each one is a difference that turned out to be a
// bug here rather than a taxonomy mismatch.
func TestCorpusFindings(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		browser Browser
	}{
		{
			// A browser's own token has to win over the "Version/4.0" that the
			// Android webview it runs in contributes.
			name:    "Opera states its version in OPR",
			agent:   "Mozilla/5.0 (Linux; U; Android 9; SDA-8TAB Build/PPR1.181005.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/119.0.6045.67 Safari/537.36 OPR/72.0.2254.67482",
			browser: Browser{BrowserOpera, Version{72, 0, 2254, "67482"}},
		},
		{
			name:    "UC Browser states its version in UCBrowser",
			agent:   "Mozilla/5.0 (Linux; U; Android 14; en-US; 24076RP19G Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 UCBrowser/13.7.8.1322 Mobile Safari/537.36",
			browser: Browser{BrowserUCBrowser, Version{13, 7, 8, "1322"}},
		},
		{
			name:    "Samsung Internet reports a version at all",
			agent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/29.0 Chrome/136.0.0.0 Safari/537.36",
			browser: Browser{BrowserSamsung, Version{29, 0, 0, ""}},
		},
		{
			name:    "Yandex browser inside the search app",
			agent:   "Mozilla/5.0 (Linux; arm_64; Android 10; AirTouch PERFORMANCE 10x) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.533 YaSearchBrowser/24.17/apad BroPP/1.0 YaSearchApp/24.17/apad webOmni SA/3 Mobile Safari/537.36",
			browser: Browser{BrowserYandex, Version{24, 17, 0, ""}},
		},
		{
			// A bare "yandex" also appears in the browsers that vendor ships, so
			// it must not be read as the crawler.
			name:    "Yandex search app is not the Yandex crawler",
			agent:   "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36 YandexSearch/23.1.0.29",
			browser: Browser{BrowserYandex, Version{23, 1, 0, "29"}},
		},
		{
			// The Facebook app sometimes sends only its own field.
			name:    "Facebook without an engine token",
			agent:   "[FBAN/FB4A;FBAV/26.0.0.22.16;FBBV/6590638;FBLC/en_US;]",
			browser: Browser{BrowserFacebook, Version{26, 0, 0, "22.16"}},
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

// The Yandex crawler must still be recognised as a bot.
func TestYandexBotStillDetected(t *testing.T) {
	for _, agent := range []string{
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		"Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots) YandexBot/3.0",
		"YaDirectFetcher/1.0 (Dyatel; +http://yandex.com/bots)",
	} {
		got := Parse(agent)
		if got.Browser.Name != BrowserYandexBot {
			t.Errorf("Browser.Name = %v, want BrowserYandexBot for %q", got.Browser.Name, agent)
		}
		if !got.IsBot() {
			t.Errorf("IsBot() = false, want true for %q", agent)
		}
	}
}

func TestIsFrozenOSVersion(t *testing.T) {
	tests := []struct {
		name  string
		agent string
		want  bool
	}{
		{
			name:  "macOS Safari 26 cannot be Catalina",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Safari/605.1.15",
			want:  true,
		},
		{
			name:  "macOS Chrome past the Catalina cutoff",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:  true,
		},
		{
			name:  "Firefox cap can still be genuine Catalina",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0",
			want:  false,
		},
		{
			name:  "similar Chrome platform token is not the exact cap",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_70) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:  false,
		},
		{
			name:  "iPad in desktop mode shares the Mac token",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want:  true,
		},
		{
			name:  "Safari 15 on Catalina may be genuine",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Safari/605.1.15",
			want:  false,
		},
		{
			name:  "a real older macOS is not frozen",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Safari/605.1.15",
			want:  false,
		},
		{
			// The token is frozen whatever the client. Parse(ua, true) can
			// recover the release from full Safari's version, which does not make
			// the token itself any less of a constant.
			name:  "iOS Safari 26 states a frozen token",
			agent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want:  true,
		},
		{
			name:  "iOS WKWebView states a frozen token",
			agent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
			want:  true,
		},
		{
			name:  "a genuine iOS 18.5 is not frozen",
			agent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Mobile/15E148 Safari/604.1",
			want:  false,
		},
		{
			name:  "Windows is unaffected",
			agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsFrozenOSVersion(test.agent); got != test.want {
				t.Fatalf("IsFrozenOSVersion() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMinimumMacOSVersion(t *testing.T) {
	tests := []struct {
		name  string
		agent string
		want  Version
	}{
		{
			name:  "Safari 26 rules out anything before Sonoma",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Safari/605.1.15",
			want:  Version{Major: 14},
		},
		{
			name:  "Safari 18 rules out anything before Ventura",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15",
			want:  Version{Major: 13},
		},
		{
			name:  "Safari 17 rules out anything before Monterey",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
			want:  Version{Major: 12},
		},
		{
			name:  "Safari 16 rules out anything before Big Sur",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
			want:  Version{Major: 11},
		},
		{
			name:  "a future Safari keeps the newest known floor",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/27.0 Safari/605.1.15",
			want:  Version{Major: 14},
		},
		{
			name:  "Safari 15 offers no useful floor",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Safari/605.1.15",
			want:  Version{},
		},
		{
			name:  "Chrome carries no Safari version to work from",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:  Version{},
		},
		{
			name:  "iOS is unaffected, the exact version is known there",
			agent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want:  Version{},
		},
		{
			name:  "visible iPad desktop marker is not a macOS floor",
			agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want:  Version{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MinimumMacOSVersion(test.agent); got != test.want {
				t.Fatalf("MinimumMacOSVersion() = %#v, want %#v", got, test.want)
			}
		})
	}
}

// The frozen marker keeps the stated numbers usable while saying they are a
// constant. It rides in Extra, which OS versions never use otherwise, and it does
// not disturb version comparisons.
func TestFrozenVersionMarker(t *testing.T) {
	tests := []struct {
		name   string
		agent  string
		want   Version
		frozen bool
	}{
		{
			name:   "iOS Safari 26 is resolved, so nothing is marked",
			agent:  "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want:   Version{26, 5, 0, ""},
			frozen: false,
		},
		{
			name:   "Safari version suffix does not become OS extra data",
			agent:  "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5.0.build Mobile/15E148 Safari/604.1",
			want:   Version{26, 5, 0, ""},
			frozen: false,
		},
		{
			name:   "iOS WKWebView keeps the numbers and is marked",
			agent:  "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
			want:   Version{18, 7, 0, VersionFrozen},
			frozen: true,
		},
		{
			name:   "macOS Safari 26 cannot be resolved, only marked",
			agent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Safari/605.1.15",
			want:   Version{10, 15, 7, VersionFrozen},
			frozen: true,
		},
		{
			name:   "macOS Chrome past the Catalina cutoff is marked",
			agent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:   Version{10, 15, 7, VersionFrozen},
			frozen: true,
		},
		{
			name:   "Firefox cap may still describe genuine Catalina",
			agent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0",
			want:   Version{10, 15, 0, ""},
			frozen: false,
		},
		{
			name:   "macOS Safari 15 may be genuine Catalina",
			agent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Safari/605.1.15",
			want:   Version{10, 15, 7, ""},
			frozen: false,
		},
		{
			name:   "Windows is untouched",
			agent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
			want:   Version{10, 0, 0, ""},
			frozen: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Parse(test.agent, true)
			if got.OS.Version != test.want {
				t.Fatalf("OS.Version = %#v, want %#v", got.OS.Version, test.want)
			}
			if got.OS.Version.IsFrozen() != test.frozen {
				t.Fatalf("IsFrozen() = %v, want %v", got.OS.Version.IsFrozen(), test.frozen)
			}

			// Parse never marks, and the numbers it reports are the same.
			plain := Parse(test.agent)
			if plain.OS.Version.IsFrozen() {
				t.Fatal("Parse marked a version, want no marker")
			}
			if got.OS.Version.IsFrozen() &&
				(plain.OS.Version.Major != got.OS.Version.Major || plain.OS.Version.Minor != got.OS.Version.Minor) {
				t.Fatalf("marked version %#v changed the numbers reported by Parse %#v", got.OS.Version, plain.OS.Version)
			}

			// A browser version must not inherit the marker.
			if got.Browser.Version.IsFrozen() {
				t.Fatal("Browser.Version.IsFrozen() = true, want false")
			}
		})
	}
}

// The marker must not interfere with ordering.
func TestFrozenVersionComparison(t *testing.T) {
	frozen := Version{18, 7, 0, VersionFrozen}
	if !frozen.Less(Version{26, 0, 0, ""}) {
		t.Error("a marked 18.7 should sort before 26.0")
	}
	older := Version{18, 6, 0, ""}
	if !older.Less(frozen) {
		t.Error("18.6 should sort before a marked 18.7")
	}
}
