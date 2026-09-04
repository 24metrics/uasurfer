package main

import "github.com/24metrics/uasurfer"

// uapCoreBrowsers and uapCoreSystems translate uap-core's families into this
// package's constants, following the same rules as the Matomo maps: only labels
// with a counterpart here are listed, and a label may accept more than one
// constant where the two projects group differently.

var uapCoreBrowsers = map[string][]uasurfer.BrowserName{
	"Chrome":                {uasurfer.BrowserChrome},
	"Chrome Frame":          {uasurfer.BrowserChrome},
	"Chrome Mobile":         {uasurfer.BrowserChrome},
	"Chrome Mobile iOS":     {uasurfer.BrowserChrome},
	"Chrome Mobile WebView": {uasurfer.BrowserChrome},
	"Chromium":              {uasurfer.BrowserChrome},
	"HeadlessChrome":        {uasurfer.BrowserChrome},

	"Safari":        {uasurfer.BrowserSafari},
	"Mobile Safari": {uasurfer.BrowserSafari},

	"Edge":            {uasurfer.BrowserIE},
	"Edge Mobile":     {uasurfer.BrowserIE},
	"IE":              {uasurfer.BrowserIE},
	"IE Large Screen": {uasurfer.BrowserIE},
	"IE Mobile":       {uasurfer.BrowserIE},

	"Firefox":             {uasurfer.BrowserFirefox},
	"Firefox (Shiretoko)": {uasurfer.BrowserFirefox},
	"Firefox Beta":        {uasurfer.BrowserFirefox},
	"Firefox Mobile":      {uasurfer.BrowserFirefox},
	"Firefox iOS":         {uasurfer.BrowserFirefox},
	"IceCat":              {uasurfer.BrowserFirefox},

	"Opera":        {uasurfer.BrowserOpera},
	"Opera Mini":   {uasurfer.BrowserOpera},
	"Opera Mobile": {uasurfer.BrowserOpera},
	"Opera Tablet": {uasurfer.BrowserOpera},

	"UC Browser":        {uasurfer.BrowserUCBrowser},
	"Amazon Silk":       {uasurfer.BrowserSilk},
	"QQ Browser":        {uasurfer.BrowserQQ},
	"QQ Browser Mobile": {uasurfer.BrowserQQ},
	"Samsung Internet":  {uasurfer.BrowserSamsung},
	"Yandex Browser":    {uasurfer.BrowserYandex},
	"Coc Coc":           {uasurfer.BrowserCocCoc},
	"Spotify":           {uasurfer.BrowserSpotify},
	"Maxthon":           {uasurfer.BrowserMaxthon},

	"Nokia Browser":     {uasurfer.BrowserNokia},
	"Nokia OSS Browser": {uasurfer.BrowserNokia},

	"Android": {uasurfer.BrowserAndroid},

	"BlackBerry WebKit": {uasurfer.BrowserBlackberry},

	"Facebook":          {uasurfer.BrowserFacebook},
	"Instagram":         {uasurfer.BrowserInstagram},
	"WeChat Browser":    {uasurfer.BrowserWeChat},
	"TikTok":            {uasurfer.BrowserTikTok},
	"Snapchat":          {uasurfer.BrowserSnapchat},
	"LINE":              {uasurfer.BrowserLine},
	"DuckDuckGo":        {uasurfer.BrowserDuckDuckGo},
	"DuckDuckGo Mobile": {uasurfer.BrowserDuckDuckGo},

	"Googlebot":              {uasurfer.BrowserGoogleBot},
	"Googlebot-Image":        {uasurfer.BrowserGoogleBot},
	"Googlebot-Mobile":       {uasurfer.BrowserGoogleBot},
	"Googlebot-News":         {uasurfer.BrowserGoogleBot},
	"Googlebot-Video":        {uasurfer.BrowserGoogleBot},
	"Googlebot-richsnippets": {uasurfer.BrowserGoogleBot},
	"bingbot":                {uasurfer.BrowserBingBot},
	"BingPreview":            {uasurfer.BrowserBingBot},
	"BaiduImagespider":       {uasurfer.BrowserBaiduBot},
	"BaiDuSpider":            {uasurfer.BrowserBaiduBot},
	"Baiduspider":            {uasurfer.BrowserBaiduBot},
	"baiduspider":            {uasurfer.BrowserBaiduBot},
	"Baiduspider-cpro":       {uasurfer.BrowserBaiduBot},
	"Baiduspider-image":      {uasurfer.BrowserBaiduBot},
	"Baiduspider-testbranch": {uasurfer.BrowserBaiduBot},
	"LinkedInBot":            {uasurfer.BrowserLinkedInBot},
	"msnbot":                 {uasurfer.BrowserMsnBot},
	"msnbot-media":           {uasurfer.BrowserMsnBot},
	"PingdomBot":             {uasurfer.BrowserPingdomBot},
	"Twitterbot":             {uasurfer.BrowserTwitterBot},
	"YandexBot":              {uasurfer.BrowserYandexBot},
	"Yahoo! Slurp":           {uasurfer.BrowserYahooBot},
	"PhantomJS":              {uasurfer.BrowserBot},
}

var uapCoreSystems = map[string][]uasurfer.OSName{
	"iOS":       {uasurfer.OSiOS},
	"Mac OS X":  {uasurfer.OSMacOSX},
	"Windows":   {uasurfer.OSWindows},
	"Chrome OS": {uasurfer.OSChromeOS},

	// A Fire device runs Android. uap-core reports it either way depending on
	// the fixture, this package always reports Kindle.
	"Android": {uasurfer.OSAndroid, uasurfer.OSKindle},
	"Kindle":  {uasurfer.OSKindle, uasurfer.OSAndroid},

	"Windows Mobile": {uasurfer.OSWindowsPhone},
	"Windows Phone":  {uasurfer.OSWindowsPhone},

	"BlackBerry OS":        {uasurfer.OSBlackberry},
	"BlackBerry Tablet OS": {uasurfer.OSBlackberry},

	"Web0S": {uasurfer.OSWebOS},
	"webOS": {uasurfer.OSWebOS},

	"Arch Linux": {uasurfer.OSLinux},
	"CentOS":     {uasurfer.OSLinux},
	"Debian":     {uasurfer.OSLinux},
	"Fedora":     {uasurfer.OSLinux},
	"Gentoo":     {uasurfer.OSLinux},
	"Linux":      {uasurfer.OSLinux},
	"Linux Mint": {uasurfer.OSLinux},
	"Mandriva":   {uasurfer.OSLinux},
	"Red Hat":    {uasurfer.OSLinux},
	"SUSE":       {uasurfer.OSLinux},
	"Slackware":  {uasurfer.OSLinux},
	"Ubuntu":     {uasurfer.OSLinux},
	"openSUSE":   {uasurfer.OSLinux},
}
