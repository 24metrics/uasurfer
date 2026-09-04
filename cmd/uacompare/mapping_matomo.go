package main

import "github.com/24metrics/uasurfer"

// The maps below translate Matomo's vocabulary into this package's constants.
// Matomo names every product it knows, this package deliberately does not, so
// only labels with a counterpart here are listed. Anything missing is reported
// as out of scope instead of counting as a miss.
//
// Some entries accept more than one constant. That is not a way to inflate the
// result: it covers the places where the two projects slice the same reality
// differently, for example a Fire tablet, which Matomo reports as Fire OS on
// Android while this package has a dedicated Kindle constant.

var matomoClients = map[string][]uasurfer.BrowserName{
	// Chromium and the names Chrome carries per platform.
	"Chrome":            {uasurfer.BrowserChrome},
	"Chrome Mobile":     {uasurfer.BrowserChrome},
	"Chrome Mobile iOS": {uasurfer.BrowserChrome},
	"Chrome Webview":    {uasurfer.BrowserChrome},
	"Chromium":          {uasurfer.BrowserChrome},
	"Headless Chrome":   {uasurfer.BrowserChrome},

	// Apple.
	"Safari":        {uasurfer.BrowserSafari},
	"Mobile Safari": {uasurfer.BrowserSafari},

	// Microsoft, which this package reports as BrowserIE including Edge.
	"Internet Explorer": {uasurfer.BrowserIE},
	"IE Mobile":         {uasurfer.BrowserIE},
	"Microsoft Edge":    {uasurfer.BrowserIE},

	// Mozilla and the rebuilds this package folds into Firefox.
	"Firefox":            {uasurfer.BrowserFirefox},
	"Firefox Mobile":     {uasurfer.BrowserFirefox},
	"Firefox Mobile iOS": {uasurfer.BrowserFirefox},
	"Firefox Beta":       {uasurfer.BrowserFirefox},
	"Iceweasel":          {uasurfer.BrowserFirefox},
	"IceCat":             {uasurfer.BrowserFirefox},
	"SeaMonkey":          {uasurfer.BrowserFirefox},

	"Opera":        {uasurfer.BrowserOpera},
	"Opera Mini":   {uasurfer.BrowserOpera},
	"Opera Mobile": {uasurfer.BrowserOpera},
	"Opera Next":   {uasurfer.BrowserOpera},

	"UC Browser":      {uasurfer.BrowserUCBrowser},
	"UC Browser Mini": {uasurfer.BrowserUCBrowser},
	"UC Browser HD":   {uasurfer.BrowserUCBrowser},

	"Amazon Silk": {uasurfer.BrowserSilk},
	"Mobile Silk": {uasurfer.BrowserSilk},

	"QQ Browser":        {uasurfer.BrowserQQ},
	"QQ Browser Mini":   {uasurfer.BrowserQQ},
	"QQ Browser Mobile": {uasurfer.BrowserQQ},

	"Samsung Browser":     {uasurfer.BrowserSamsung},
	"Yandex Browser":      {uasurfer.BrowserYandex},
	"Yandex Browser Lite": {uasurfer.BrowserYandex},
	"Coc Coc":             {uasurfer.BrowserCocCoc},
	"Maxthon":             {uasurfer.BrowserMaxthon},
	"Sogou Explorer":      {uasurfer.BrowserSogouExplorer},
	"Nintendo Browser":    {uasurfer.BrowserNintendo},
	"Spotify":             {uasurfer.BrowserSpotify},

	"Android Browser": {uasurfer.BrowserAndroid},

	"NetFront":      {uasurfer.BrowserNetFront},
	"NetFront Life": {uasurfer.BrowserNetFront},

	"BlackBerry Browser": {uasurfer.BrowserBlackberry},

	"Nokia Browser":     {uasurfer.BrowserNokia},
	"Nokia OSS Browser": {uasurfer.BrowserNokia},
	"Nokia Ovi Browser": {uasurfer.BrowserNokia},

	// In-app browsers. Matomo separates a vendor's apps, this package reports
	// one constant per vendor.
	"Facebook":                   {uasurfer.BrowserFacebook},
	"Facebook Messenger":         {uasurfer.BrowserFacebook},
	"Facebook Lite":              {uasurfer.BrowserFacebook},
	"Instagram":                  {uasurfer.BrowserInstagram},
	"Instagram App":              {uasurfer.BrowserInstagram},
	"WeChat":                     {uasurfer.BrowserWeChat},
	"WeChat Share Extension":     {uasurfer.BrowserWeChat},
	"TikTok":                     {uasurfer.BrowserTikTok},
	"Snapchat":                   {uasurfer.BrowserSnapchat},
	"LINE":                       {uasurfer.BrowserLine},
	"DuckDuckGo Privacy Browser": {uasurfer.BrowserDuckDuckGo},
}

var matomoSystems = map[string][]uasurfer.OSName{
	"Android":     {uasurfer.OSAndroid},
	"Android TV":  {uasurfer.OSAndroid},
	"Lineage OS":  {uasurfer.OSAndroid},
	"iOS":         {uasurfer.OSiOS},
	"Mac":         {uasurfer.OSMacOSX},
	"Chrome OS":   {uasurfer.OSChromeOS},
	"Chromium OS": {uasurfer.OSChromeOS},
	"webOS":       {uasurfer.OSWebOS},
	"PlayStation": {uasurfer.OSPlaystation},
	"Xbox":        {uasurfer.OSXbox},
	"Nintendo":    {uasurfer.OSNintendo},

	// A console running Windows is reported by its system here and by the
	// platform underneath in Matomo, so both readings are accepted.
	"Windows":        {uasurfer.OSWindows, uasurfer.OSXbox},
	"Windows RT":     {uasurfer.OSWindows, uasurfer.OSXbox},
	"Windows Phone":  {uasurfer.OSWindowsPhone, uasurfer.OSXbox},
	"Windows Mobile": {uasurfer.OSWindowsPhone},
	"Windows CE":     {uasurfer.OSWindowsPhone},

	// A Fire device runs Android; this package reports it as Kindle.
	"Fire OS": {uasurfer.OSKindle, uasurfer.OSAndroid},

	"BlackBerry OS":        {uasurfer.OSBlackberry},
	"BlackBerry Tablet OS": {uasurfer.OSBlackberry},

	// Every distribution collapses into one constant here.
	"GNU/Linux":  {uasurfer.OSLinux},
	"Arch Linux": {uasurfer.OSLinux},
	"CentOS":     {uasurfer.OSLinux},
	"Debian":     {uasurfer.OSLinux},
	"Fedora":     {uasurfer.OSLinux},
	"Gentoo":     {uasurfer.OSLinux},
	"Knoppix":    {uasurfer.OSLinux},
	"Kubuntu":    {uasurfer.OSLinux},
	"Linux Mint": {uasurfer.OSLinux},
	"Mandriva":   {uasurfer.OSLinux},
	"Raspbian":   {uasurfer.OSLinux},
	"Red Hat":    {uasurfer.OSLinux},
	"SUSE":       {uasurfer.OSLinux},
	"Slackware":  {uasurfer.OSLinux},
	"Ubuntu":     {uasurfer.OSLinux},
	"openSUSE":   {uasurfer.OSLinux},
}

var matomoDevices = map[string][]uasurfer.DeviceType{
	"desktop":       {uasurfer.DeviceComputer},
	"smartphone":    {uasurfer.DevicePhone},
	"feature phone": {uasurfer.DevicePhone},
	"tablet":        {uasurfer.DeviceTablet},
	"console":       {uasurfer.DeviceConsole},
	"wearable":      {uasurfer.DeviceWearable},

	// A phablet is a large phone; this package has no such class and may read it
	// either way depending on the tablet hints in the user agent.
	"phablet": {uasurfer.DevicePhone, uasurfer.DeviceTablet},

	// Matomo separates a television from the box in front of it.
	"tv":            {uasurfer.DeviceTV, uasurfer.DeviceMediaHub},
	"smart display": {uasurfer.DeviceTV, uasurfer.DeviceMediaHub},
}
