[![Build Status](https://travis-ci.org/avct/uasurfer.svg?branch=master)](https://travis-ci.org/avct/uasurfer)  [![GoDoc](https://godoc.org/github.com/avct/uasurfer?status.svg)](https://godoc.org/github.com/avct/uasurfer)  [![Go Report Card](https://goreportcard.com/badge/github.com/avct/uasurfer)](https://goreportcard.com/report/github.com/avct/uasurfer)

# uasurfer

![uasurfer-100px](https://cloud.githubusercontent.com/assets/597902/16172506/9debc136-357a-11e6-90fb-c7c46f50dff0.png)

**User Agent Surfer** (uasurfer) is a lightweight Golang package that parses and abstracts [HTTP User-Agent strings](https://en.wikipedia.org/wiki/User_agent) with particular attention to device type.

The following information is returned by uasurfer from a raw HTTP User-Agent string:

| Name           | Example | Coverage in 192,792 parses |
|----------------|---------|--------------------------------|
| Browser name    | `chrome` | 99.85%                         |
| Browser version | `53` | 99.17%                         |
| Platform       | `ipad`  | 99.97%                         |
| OS name         | `ios`  | 99.96%                         |
| OS version      | `10`   | 98.81%                         |
| Device type    |  `tablet` | 99.98%                         |

Layout engine, browser language, and other esoteric attributes are not parsed.

Coverage is estimated from a random sample of real UA strings collected across thousands of sources in US and EU mid-2016.

## Usage

### Parse(ua string) Function

The `Parse()` function accepts a user agent `string` and returns a `UserAgent` struct with named constants and separate major, minor, and patch version fields. Called with one argument it reports what the user agent states and applies no heuristics.

That matters because browsers now cap the OS version they report. Current WebKit freezes the iOS-family token at `18_7` (early iOS 26 builds used `18_6` or `18_6_2`), and on the desktop the token stays at `Intel Mac OS X 10_15_7`, which iPad Safari sends too. `Parse()` passes those values through, as every other parser does. Four additions let you deal with them:

* `IsFrozenOSVersion()` tells a capped value from a genuine one.
* `MinimumMacOSVersion()` derives a lower bound from the Safari version, the only usable signal left on the Mac.
* `ParseOSVersionDetails()` returns the reported value, frozen status and optional floor together with one parse.
* The optional second argument to `Parse()` resolves the iOS-family version where the Safari version identifies the release, and marks it as frozen where it does not. It never returns an unknown version in place of a capped one.

### Resolving frozen OS versions

```go
func Parse(ua string, resolveFrozenOSVersion ...bool) *UserAgent
func ParseUserAgent(ua string, dest *UserAgent, resolveFrozenOSVersion ...bool)
```

```go
ua := uasurfer.Parse(rawUA)        // as stated
ua := uasurfer.Parse(rawUA, false) // identical
ua := uasurfer.Parse(rawUA, true)  // resolved where possible, marked otherwise
```

Existing direct `Parse(rawUA)` calls remain source-compatible. The function's Go
type is now `func(string, ...bool) *UserAgent`, however, so code that stores it in
a `func(string) *UserAgent` variable must wrap it:

```go
parse := func(rawUA string) *uasurfer.UserAgent { return uasurfer.Parse(rawUA) }
```

Mobile Safari ships with iOS, so its version identifies the release. A capped
version is therefore either replaced by the release it can be recovered from, or
kept and marked. The stated numbers are never discarded:

| user agent | `Parse()` | `Parse(ua, true)` |
| --- | --- | --- |
| iPhone, `18_7`, Safari 26.5 | `{18, 7, 0}` | `{26, 5, 0}` |
| iPhone, `18_7`, bare WKWebView | `{18, 7, 0}` | `{18, 7, 0, "frozen"}` |
| iPhone, `18_7`, Safari 26 plus an app marker | `{18, 7, 0}` | `{18, 7, 0, "frozen"}` |
| iPhone, `18_7`, `FxiOS` or `CriOS` | `{18, 7, 0}` | `{18, 7, 0, "frozen"}` |
| iPhone, `18_5`, Safari 18.5 | `{18, 5, 0}` | `{18, 5, 0}` |
| macOS, `10_15_7`, Safari 26.5 | `{10, 15, 7}` | `{10, 15, 7, "frozen"}` |
| macOS, `10_15_7`, Safari 15.6 | `{10, 15, 7}` | `{10, 15, 7}` |

Only full Safari on iOS gets its release recovered: a WKWebView, an in-app browser
or a third-party browser carries no hint about the system it runs on, and on macOS
nothing can be recovered at all, so those are marked instead. Use
`MinimumMacOSVersion()` for the lower bound that is still available there.

### The frozen marker

A marked version carries `VersionFrozen` in `Extra` and answers `true` to
`IsFrozen()`. The numbers stay the ones stated in the user agent, so they remain
usable, and `Less()` ignores the marker, so ordering and comparisons are unaffected:

```go
ua := uasurfer.Parse(rawUA, true)

switch {
case ua.OS.Version.IsFrozen():
    // Numbers are a constant, not a measurement. Bucket them separately,
    // or fall back to a floor via MinimumMacOSVersion.
case ua.OS.Version == (uasurfer.Version{}):
    // No version in the user agent at all.
default:
    // Usable version.
}
```

The parser assigns this marker only to `OS.Version`. For a browser version
`Extra` keeps its existing meaning, the components beyond the patch level, as in
Chrome 30.0.1599.101 leaving `"101"`. The macOS floor is deliberately **not**
stored in `Extra`: `MinimumMacOSVersion(rawUA)` returns it separately as a
`Version` and does not modify the result of `Parse`. For example, Safari 26 can
produce `OS.Version == Version{10, 15, 7, "frozen"}` while
`MinimumMacOSVersion(rawUA) == Version{Major: 14}`.

Because `IsFrozen()` is a method on the shared `Version` type, call it on the OS
version returned by `Parse(ua, true)`; an arbitrary browser token whose fourth
component is literally `frozen` has the same `Extra` value. `Parse()` without the
flag never assigns the marker.

Keeping this out of `Parse()` has a second benefit: a comparison against an
external corpus stays a comparison of plain user agent readings. See
`cmd/uacompare` below.

### ParseOSVersionDetails(ua string) OSVersionDetails

Returns the reported OS version, frozen status and optional macOS floor together,
while parsing the user agent only once:

```go
details := uasurfer.ParseOSVersionDetails(rawUA)

// Safari 26 on a Mac:
// details.Reported == uasurfer.Version{Major: 10, Minor: 15, Patch: 7}
// details.Minimum  == uasurfer.Version{Major: 14}
// details.Frozen   == true
```

`Reported` is deliberately the plain value from `Parse(rawUA).OS.Version`. It
never contains `VersionFrozen`, and `Minimum` is not encoded in its `Extra`
field. For a frozen Full Safari UA on iOS, for example, `Reported` remains
`{18, 7, 0}` and `Frozen` is true. Use `Parse(rawUA, true)` separately if the
resolved `{26, 5, 0}` representation is wanted.

The fields have exactly the same meaning as the individual APIs:

```go
details.Reported == uasurfer.Parse(rawUA).OS.Version
details.Minimum  == uasurfer.MinimumMacOSVersion(rawUA)
details.Frozen   == uasurfer.IsFrozenOSVersion(rawUA)
```

### IsFrozenOSVersion(ua string) bool

Reports whether the `OS.Version` returned by `Parse()` is a constant the browser
sends for every release. Use it to decide whether the version can be used at all:

```go
result := uasurfer.Parse(rawUA)
if uasurfer.IsFrozenOSVersion(rawUA) {
    // The reported OS.Version is a placeholder; treat it as unknown.
}
```

On the Mac it is true only where the value provably cannot be genuine. On the iOS
family every frozen token counts, whatever client sent it, because WebKit hardcodes
it for all of them:

| user agent | `Parse()` `OS.Version` | frozen |
| --- | --- | --- |
| macOS, `10_15_7`, Safari 26.5 | `{10, 15, 7}` | yes, Catalina supports no Safari beyond 15.6 |
| macOS, `10_15_7`, Chrome 136 | `{10, 15, 7}` | yes, Catalina supports no Chrome beyond 128 |
| macOS, `10.15`, Firefox 135 | `{10, 15, 0}` | no, the cap exists but Firefox can still run on Catalina |
| iPad in desktop mode, `10_15_7`, Safari 26.5 | `{10, 15, 7}` | yes, same synthetic token as a Mac |
| macOS, `10_15_7`, Safari 15.6 | `{10, 15, 7}` | no, could be a real Catalina machine |
| macOS, `10_14_6`, Safari 13.1 | `{10, 14, 6}` | no, genuine value |
| iPhone, `18_7`, Safari 26.5 | `{18, 7, 0}` | yes, `Parse(ua, true)` recovers `{26, 5, 0}` |
| iPhone, `18_7`, bare WKWebView | `{18, 7, 0}` | yes, `Parse(ua, true)` marks it instead |
| iPhone, `18_5`, Safari 18.5 | `{18, 5, 0}` | no, genuine value |

The answer describes the user agent, not a particular call, so it stays true for an
iOS agent whose version `Parse` can recover. Note also that the frozen iOS
values are real releases as well: a device genuinely running 18.7 is reported as
frozen, because the two cannot be told apart.

A false result means the version is genuine or not provably capped, not that it is
guaranteed to be accurate.

### MinimumMacOSVersion(ua string) Version

Returns the oldest macOS release that could be behind a user agent, or the zero
`Version` when no bound can be derived. Where the capped token says nothing,
Safari's own version still rules out older systems:

```go
if floor := uasurfer.MinimumMacOSVersion(rawUA); floor.Major > 0 {
    // If this is a Mac, it runs at least this macOS version.
}
```

An iPad requesting a desktop site can send a user agent indistinguishable from a
Mac. When `Mobile/` exposes the iPad origin this function returns zero; without
such a hint the caller must already know that the request came from a Mac before
interpreting the result as an OS floor.

| Safari | available for | floor |
| --- | --- | --- |
| 26 | macOS 26, Sequoia, Sonoma | `{14, 0, 0}` |
| 18 | macOS 15, Sonoma, Ventura | `{13, 0, 0}` |
| 17 | Sonoma, Ventura, Monterey | `{12, 0, 0}` |
| 16 | Ventura, Monterey, Big Sur | `{11, 0, 0}` |

This is a floor, not a guess. Deriving an exact version by subtracting a constant
from the Safari version is wrong twice over: Safari also ships for the two
previous macOS releases, so one Safari version spans three systems, and the
year-based renumbering in 2025 broke any fixed offset — Safari 26 belongs to
macOS 26, not to macOS 23. For the same reason there is no formula in the code,
only a table taken from Apple's release notes, and a new Safari generation needs a
new entry. Other browsers cap the token without offering a usable version of their
own, so they return the zero `Version`.

The only way to obtain an exact macOS version is the `Sec-CH-UA-Platform-Version`
client hint, which Chromium sends and Safari and Firefox do not; this package
parses user agents only.

```go
// Define a user agent string
myUA := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36"

ua := uasurfer.Parse(myUA)
```

where example UserAgent is:
```
{
    Browser {
        BrowserName: BrowserChrome,
        Version: {
            Major: 45,
            Minor: 0,
            Patch: 2454,
        },
    },
    OS {
        Platform: PlatformMac,
        Name: OSMacOSX,
        Version: {
            Major: 10,
            Minor: 10,
            Patch: 5,
        },
    },
    DeviceType: DeviceComputer,
}
```

**Usage note:** There are some OSes that do not return a version, see docs below. Linux is typically not reported with a specific Linux distro name or version.

#### Browser Name
* `BrowserChrome` - Google [Chrome](https://en.wikipedia.org/wiki/Google_Chrome), [Chromium](https://en.wikipedia.org/wiki/Chromium_(web_browser))
* `BrowserSafari` - Apple [Safari](https://en.wikipedia.org/wiki/Safari_(web_browser))
* `BrowserIE` - Microsoft [Internet Explorer](https://en.wikipedia.org/wiki/Internet_Explorer), [Edge](https://en.wikipedia.org/wiki/Microsoft_Edge)
* `BrowserFirefox` - Mozilla [Firefox](https://en.wikipedia.org/wiki/Firefox), GNU [IceCat](https://en.wikipedia.org/wiki/GNU_IceCat), [Iceweasel](https://en.wikipedia.org/wiki/Mozilla_Corporation_software_rebranded_by_the_Debian_project#Iceweasel), [Seamonkey](https://en.wikipedia.org/wiki/SeaMonkey)
* `BrowserAndroid` - Android [WebView](https://developer.chrome.com/multidevice/webview/overview) (Android OS <4.4 only)
* `BrowserOpera` - [Opera](https://en.wikipedia.org/wiki/Opera_(web_browser))
* `BrowserUCBrowser` - [UC Browser](https://en.wikipedia.org/wiki/UC_Browser)
* `BrowserSilk` - Amazon [Silk](https://en.wikipedia.org/wiki/Amazon_Silk)
* `BrowserQQ` - Tencent [QQ](https://en.wikipedia.org/wiki/Tencent_QQ)
* `BrowserSpotify` - [Spotify](https://en.wikipedia.org/wiki/Spotify#Clients) desktop client
* `BrowserBlackberry` - RIM [BlackBerry](https://en.wikipedia.org/wiki/BlackBerry)
* `BrowserYandex` - [Yandex](https://en.wikipedia.org/wiki/Yandex_Browser)
* `BrowserNintendo` - [Nintendo DS(i) Browser](https://en.wikipedia.org/wiki/Nintendo_DS_%26_DSi_Browser)
* `BrowserSamsung` - [Samsung Internet](https://en.wikipedia.org/wiki/Samsung_Internet_for_Android)
* `BrowserCocCoc`- [Cốc Cốc](https://en.wikipedia.org/wiki/C%E1%BB%91c_C%E1%BB%91c)
* `BrowserFacebook` - Facebook in-app browser
* `BrowserInstagram` - Instagram in-app browser
* `BrowserWeChat` - WeChat in-app browser
* `BrowserTikTok` - TikTok in-app browser
* `BrowserSnapchat` - Snapchat in-app browser
* `BrowserLine` - Line in-app browser
* `BrowserDuckDuckGo` - DuckDuckGo browser
* `BrowserUnknown` - Unknown

#### Browser Version

Browser version returns a `Version` with major, minor, patch and optional extra
components from the User-Agent string. For example, Chrome 45.0.23423 returns a
major version of `45`. The numeric fields support comparisons such as "Chrome
version > 23".

An unknown version is represented by `Version{}`.

#### Platform
* `PlatformWindows` - Microsoft Windows
* `PlatformMac` - Apple Macintosh
* `PlatformLinux` - Linux, including Android and other OSes
* `PlatformiPad` - Apple iPad
* `PlatformiPhone` - Apple iPhone
* `PlatformBlackberry` - RIM Blackberry
* `PlatformWindowsPhone` Microsoft Windows Phone & Mobile
* `PlatformKindle` - Amazon Kindle & Kindle Fire
* `PlatformPlaystation` - Sony Playstation, Vita, PSP
* `PlatformXbox` - Microsoft Xbox
* `PlatformNintendo` - Nintendo DS, Wii, etc.
* `PlatformUnknown` - Unknown

#### OS Name
* `OSWindows`
* `OSMacOSX` - includes "macOS Sierra"
* `OSiOS`
* `OSAndroid`
* `OSChromeOS`
* `OSWebOS`
* `OSLinux`
* `OSPlaystation`
* `OSXbox`
* `OSNintendo`
* `OSUnknown`

#### OS Version

macOS versions use Apple's numeric release number, for example 10.15 for
Catalina, 11 for Big Sur, 14 for Sonoma, and 26 for Tahoe. Windows versions use
the NT version. `Version{}` indicates that a version is unknown or was not
evaluated.
Versions can be compared using `Less` function: `if ver1.Less(ver2) {}`

Here are some examples across the platform, os.name, and os.version:

* For Windows XP (Windows NT 5.1), "`PlatformWindows`" is the platform, "`OSWindows`" is the name, and `{5, 1, 0}` the version.
* For OS X 10.5.1, "`PlatformMac`" is the platform, "`OSMacOSX`" the name, and `{10, 5, 1}` the version.
* For Android 5.1, "`PlatformLinux`" is the platform, "`OSAndroid`" is the name, and `{5, 1, 0}` the version.
* For iOS 5.1, "`PlatformiPhone`" or "`PlatformiPad`" is the platform, "`OSiOS`" is the name, and `{5, 1, 0}` the version.

###### Windows Version Guide

* Windows 10 - `{10, 0, 0}`
* Windows 8.1 - `{6, 3, 0}`
* Windows 8 - `{6, 2, 0}`
* Windows 7 - `{6, 1, 0}`
* Windows Vista - `{6, 0, 0}`
* Windows XP - `{5, 1, 0}` or `{5, 2, 0}`
* Windows 2000 - `{5, 0, 0}`

Windows 95, 98, and ME represent 0.01% of traffic worldwide and are not available through this package at this time.

#### DeviceType
DeviceType is typically quite accurate, though determining between phones and tablets on Android is not always possible due to how some vendors design their UA strings. A mobile Android device without tablet indicator defaults to being classified as a phone. DeviceTV supports major brands such as Philips, Sharp, Vizio and steaming boxes such as Apple, Google, Roku, Amazon.

* `DeviceComputer`
* `DevicePhone`
* `DeviceTablet`
* `DeviceTV`
* `DeviceConsole`
* `DeviceWearable`
* `DeviceUnknown`

## Example Combinations of Attributes
* Surface RT -> `OSWindows8`, `DeviceTablet`, OSVersion >= `6`
* Android Tablet -> `OSAndroid`, `DeviceTablet`
* Microsoft Edge -> `BrowserIE`, BrowserVersion >= `12.0.0`

## Measuring hit quality against a reference

`cmd/uacompare` runs this parser against an external corpus and reports a hit
rate per dimension. Two sources are supported:

| source | project | license | cases |
| --- | --- | --- | --- |
| `matomo` (default) | [matomo-org/device-detector](https://github.com/matomo-org/device-detector) | LGPL-3.0-or-later | ~36,000 |
| `uap-core` | [ua-parser/uap-core](https://github.com/ua-parser/uap-core) | Apache-2.0 | ~2,000 |

```bash
go run ./cmd/uacompare -fetch     # download the fixtures once
go run ./cmd/uacompare            # compare, with the most frequent differences
go run ./cmd/uacompare -show 0    # summary only
go run ./cmd/uacompare -source uap-core -fetch
go run ./cmd/uacompare -strict    # non-zero exit while mismatches remain
```

The fixtures are downloaded on demand and are git-ignored. That matters for the
Matomo corpus, whose data carries the LGPL: fetching it locally to measure
against is use, not redistribution.

Current result against Matomo:

| dimension | rate |
| --- | --- |
| Browser name | 98.8% |
| Browser version | 99.6% |
| OS name | 97.8% |
| OS version | 97.8% |
| Device type | 75.6% |

### Reading the numbers

The references name every product they know, this package deliberately exposes a
small set of constants, so a raw comparison would mostly measure vocabulary. Two
mechanisms keep the result meaningful:

* `cmd/uacompare/mapping_*.go` translates reference labels into these constants.
  A label without a counterpart here counts as **out of scope**, not as a miss,
  and a label may accept more than one constant where the projects group
  differently, for example a Fire tablet being Android in Matomo and Kindle here.
* `testdata/deviations.tsv` declares intentional differences, such as Windows
  being reported by its NT version rather than its marketing name. These count as
  **declared** and are listed with a reason.

The comparison calls `Parse()` without the optional flag, so it measures plain user
agent readings against the reference. That is the reason resolving frozen versions
is opt-in: the references pass the frozen token through as well, so a resolved
value would show up as a mismatch against data that is not wrong, only capped.

Device type is the weakest dimension by design: Matomo resolves it from a device
model database, so it knows that a given Android model is a game console or a TV
box, which the user agent text alone does not say. Treat that rate as a lower
bound rather than as a defect count.

The uap-core rates are lower across the board because that corpus is weighted
toward legacy and exotic agents, including pre-XP Windows and Linux distribution
versions, which this package does not model. It is useful as a source of edge
cases, not as a picture of real traffic.

## To do

* Remove compiled regexp in favor of string.Contains wherever possible (lowers mem/alloc)
* Better version support on Firefox derivatives (e.g. SeaMonkey)
* Potential additional browser support:
 * "NetFront" (1% share in India)
 * "Sogou Explorer" (5% share in China)
 * "Maxthon" (1.5% share in China)
 * "Nokia"
* Potential additional OS support:
 * "Nokia" (5% share in India)
 * "Series 40" (5.5% share in India)
 * Windows 2003 Server
* iOS safari browser identification based on iOS version
* Add android version to browser identification
* old Macs
 * "opera/9.64 (macintosh; ppc mac os x; u; en) presto/2.1.1"
* old Windows
 * "mozilla/5.0 (windows nt 4.0; wow64) applewebkit/537.36 (khtml, like gecko) chrome/37.0.2049.0 safari/537.36"
