package uasurfer

import (
	"reflect"
	"testing"
)

// The iOS family reports a frozen platform token. Parse passes it through, and
// Parse(ua, true) resolves it: full Safari ships with the system, so its version
// identifies the release, while a WKWebView or a third-party browser carries no
// such hint and keeps the stated numbers with the frozen marker instead.
func TestIOS26AndBrowserTokens(t *testing.T) {
	tests := []struct {
		name         string
		ua           string
		want         UserAgent
		wantResolved UserAgent
	}{
		{
			name: "iOS 26.5 full Safari",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.5 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserSafari, Version{26, 5, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{26, 5, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{26, 5, 0, ""}},
				DevicePhone,
			},
		},
		{
			name: "ambiguous bare WKWebView",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
			want: UserAgent{
				Browser{BrowserSafari, Version{18, 7, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{18, 7, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, VersionFrozen}},
				DevicePhone,
			},
		},
		{
			name: "pre-freeze Safari keeps the reported release",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/25.6 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserSafari, Version{25, 6, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{25, 6, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, VersionFrozen}},
				DevicePhone,
			},
		},
		{
			name: "early iOS 26 frozen at 18.6",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserSafari, Version{26, 0, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 6, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{26, 0, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{26, 0, 0, ""}},
				DevicePhone,
			},
		},
		{
			name: "early iOS 26 frozen at 18.6.2",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_6_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserSafari, Version{26, 0, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 6, 2, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{26, 0, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{26, 0, 0, ""}},
				DevicePhone,
			},
		},
		{
			name: "iPadOS 26 mobile Safari",
			ua:   "Mozilla/5.0 (iPad; CPU OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.2 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserSafari, Version{26, 2, 0, ""}},
				OS{PlatformiPad, OSiPadOS, Version{18, 7, 0, ""}},
				DeviceTablet,
			},
			wantResolved: UserAgent{
				Browser{BrowserSafari, Version{26, 2, 0, ""}},
				OS{PlatformiPad, OSiPadOS, Version{26, 2, 0, ""}},
				DeviceTablet,
			},
		},
		{
			name: "FxiOS token wins over Version",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 FxiOS/130.1 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserFirefox, Version{130, 1, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserFirefox, Version{130, 1, 0, ""}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, VersionFrozen}},
				DevicePhone,
			},
		},
		{
			name: "OPiOS token wins over Version",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 OPiOS/10.2.0.93022 Mobile/15E148 Safari/604.1",
			want: UserAgent{
				Browser{BrowserOpera, Version{10, 2, 0, "93022"}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, ""}},
				DevicePhone,
			},
			wantResolved: UserAgent{
				Browser{BrowserOpera, Version{10, 2, 0, "93022"}},
				OS{PlatformiPhone, OSiOS, Version{18, 7, 0, VersionFrozen}},
				DevicePhone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Every entry point has to agree, both with and without the
			// heuristic, and passing false has to equal Parse.
			plain := []func(string) *UserAgent{
				func(ua string) *UserAgent { return Parse(ua) },
				func(ua string) *UserAgent {
					got := new(UserAgent)
					ParseUserAgent(ua, got)
					return got
				},
				func(ua string) *UserAgent { return Parse(ua, false) },
				func(ua string) *UserAgent {
					got := new(UserAgent)
					ParseUserAgent(ua, got, false)
					return got
				},
			}
			for _, parseUA := range plain {
				if got := *parseUA(test.ua); !reflect.DeepEqual(got, test.want) {
					t.Fatalf("Parse result = %+v, want %+v", got, test.want)
				}
			}

			resolved := []func(string) *UserAgent{
				func(ua string) *UserAgent { return Parse(ua, true) },
				func(ua string) *UserAgent {
					got := new(UserAgent)
					ParseUserAgent(ua, got, true)
					return got
				},
			}
			for _, parseUA := range resolved {
				if got := *parseUA(test.ua); !reflect.DeepEqual(got, test.wantResolved) {
					t.Fatalf("Parse result = %+v, want %+v", got, test.wantResolved)
				}
			}
		})
	}
}
