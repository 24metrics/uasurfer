package uasurfer

import "testing"

func TestForkCompatibility(t *testing.T) {
	tests := []struct {
		name    string
		agent   string
		browser BrowserName
		version Version
		os      OSName
		device  DeviceType
	}{
		{
			name:    "version extra",
			agent:   "Mozilla/5.0 AppleWebKit/537.36 Chrome/83.0.4103.120 Safari/537.36",
			browser: BrowserChrome,
			version: Version{83, 0, 4103, "120"},
			device:  DeviceUnknown,
		},
		{
			name:    "dalvik nexus player",
			agent:   "Dalvik/2.1.0 (Linux; U; Android 6.0.1; Nexus Player Build/MMB29T)",
			browser: BrowserAndroid,
			version: Version{2, 1, 0, ""},
			os:      OSAndroid,
			device:  DeviceConsole,
		},
		{
			name:    "x88 tv",
			agent:   "Mozilla/5.0 (Linux; Android 11; X88pro10.r1.00.6330.d4 Build/RP1A.201105.002; wv) AppleWebKit/537.36 Chrome/83.0.4103.120 Safari/537.36",
			browser: BrowserChrome,
			version: Version{83, 0, 4103, "120"},
			os:      OSAndroid,
			device:  DeviceTV,
		},
		{
			name:    "chromecast media hub",
			agent:   "Mozilla/5.0 (CrKey armv7l 1.4.15250) AppleWebKit/537.36 Chrome/31.0.1650.0 Safari/537.36",
			browser: BrowserChrome,
			version: Version{31, 0, 1650, "0"},
			device:  DeviceMediaHub,
		},
		{
			name:    "roku media hub",
			agent:   "Roku/DVP-5.2 (025.02E03197A)",
			browser: BrowserUnknown,
			version: Version{},
			os:      OSRoku,
			device:  DeviceMediaHub,
		},
		{
			name:    "xbox windows phone signature",
			agent:   "Mozilla/5.0 (Windows Phone 10.0; Android 4.2.1; Xbox; Xbox One) AppleWebKit/537.36 Chrome/46.0.2486.0 Mobile Safari/537.36 Edge/13.10586",
			browser: BrowserIE,
			version: Version{13, 10586, 0, ""},
			os:      OSXbox,
			device:  DeviceConsole,
		},
		{
			name:    "nintendo browser version",
			agent:   "Mozilla/5.0 (Nintendo WiiU) AppleWebKit/536.30 NX/3.0.4.2.12 NintendoBrowser/4.3.1.11264.US",
			browser: BrowserNintendo,
			version: Version{4, 3, 1, "11264.us"},
			os:      OSNintendo,
			device:  DeviceConsole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.agent)
			if got.Browser.Name != tt.browser || got.Browser.Version != tt.version || got.OS.Name != tt.os || got.DeviceType != tt.device {
				t.Fatalf("Parse() = %+v, want browser=%s version=%+v os=%s device=%s", got, tt.browser, tt.version, tt.os, tt.device)
			}
		})
	}
}
