// Code generated by "enumer -type=BrowserName -json -bson"; DO NOT EDIT.

package uasurfer

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strings"
)

const _BrowserNameName = "BrowserUnknownBrowserChromeBrowserIEBrowserSafariBrowserFirefoxBrowserAndroidBrowserOperaBrowserBlackberryBrowserUCBrowserBrowserSilkBrowserNokiaBrowserNetFrontBrowserQQBrowserMaxthonBrowserSogouExplorerBrowserSpotifyBrowserNintendoBrowserSamsungBrowserYandexBrowserCocCocBrowserBotBrowserAppleBotBrowserBaiduBotBrowserBingBotBrowserDuckDuckGoBotBrowserFacebookBotBrowserGoogleBotBrowserLinkedInBotBrowserMsnBotBrowserPingdomBotBrowserTwitterBotBrowserYandexBotBrowserCocCocBotBrowserYahooBot"

var _BrowserNameIndex = [...]uint16{0, 14, 27, 36, 49, 63, 77, 89, 106, 122, 133, 145, 160, 169, 183, 203, 217, 232, 246, 259, 272, 282, 297, 312, 326, 346, 364, 380, 398, 411, 428, 445, 461, 477, 492}

const _BrowserNameLowerName = "browserunknownbrowserchromebrowseriebrowsersafaribrowserfirefoxbrowserandroidbrowseroperabrowserblackberrybrowserucbrowserbrowsersilkbrowsernokiabrowsernetfrontbrowserqqbrowsermaxthonbrowsersogouexplorerbrowserspotifybrowsernintendobrowsersamsungbrowseryandexbrowsercoccocbrowserbotbrowserapplebotbrowserbaidubotbrowserbingbotbrowserduckduckgobotbrowserfacebookbotbrowsergooglebotbrowserlinkedinbotbrowsermsnbotbrowserpingdombotbrowsertwitterbotbrowseryandexbotbrowsercoccocbotbrowseryahoobot"

func (i BrowserName) String() string {
	if i < 0 || i >= BrowserName(len(_BrowserNameIndex)-1) {
		return fmt.Sprintf("BrowserName(%d)", i)
	}
	return _BrowserNameName[_BrowserNameIndex[i]:_BrowserNameIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _BrowserNameNoOp() {
	var x [1]struct{}
	_ = x[BrowserUnknown-(0)]
	_ = x[BrowserChrome-(1)]
	_ = x[BrowserIE-(2)]
	_ = x[BrowserSafari-(3)]
	_ = x[BrowserFirefox-(4)]
	_ = x[BrowserAndroid-(5)]
	_ = x[BrowserOpera-(6)]
	_ = x[BrowserBlackberry-(7)]
	_ = x[BrowserUCBrowser-(8)]
	_ = x[BrowserSilk-(9)]
	_ = x[BrowserNokia-(10)]
	_ = x[BrowserNetFront-(11)]
	_ = x[BrowserQQ-(12)]
	_ = x[BrowserMaxthon-(13)]
	_ = x[BrowserSogouExplorer-(14)]
	_ = x[BrowserSpotify-(15)]
	_ = x[BrowserNintendo-(16)]
	_ = x[BrowserSamsung-(17)]
	_ = x[BrowserYandex-(18)]
	_ = x[BrowserCocCoc-(19)]
	_ = x[BrowserBot-(20)]
	_ = x[BrowserAppleBot-(21)]
	_ = x[BrowserBaiduBot-(22)]
	_ = x[BrowserBingBot-(23)]
	_ = x[BrowserDuckDuckGoBot-(24)]
	_ = x[BrowserFacebookBot-(25)]
	_ = x[BrowserGoogleBot-(26)]
	_ = x[BrowserLinkedInBot-(27)]
	_ = x[BrowserMsnBot-(28)]
	_ = x[BrowserPingdomBot-(29)]
	_ = x[BrowserTwitterBot-(30)]
	_ = x[BrowserYandexBot-(31)]
	_ = x[BrowserCocCocBot-(32)]
	_ = x[BrowserYahooBot-(33)]
}

var _BrowserNameValues = []BrowserName{BrowserUnknown, BrowserChrome, BrowserIE, BrowserSafari, BrowserFirefox, BrowserAndroid, BrowserOpera, BrowserBlackberry, BrowserUCBrowser, BrowserSilk, BrowserNokia, BrowserNetFront, BrowserQQ, BrowserMaxthon, BrowserSogouExplorer, BrowserSpotify, BrowserNintendo, BrowserSamsung, BrowserYandex, BrowserCocCoc, BrowserBot, BrowserAppleBot, BrowserBaiduBot, BrowserBingBot, BrowserDuckDuckGoBot, BrowserFacebookBot, BrowserGoogleBot, BrowserLinkedInBot, BrowserMsnBot, BrowserPingdomBot, BrowserTwitterBot, BrowserYandexBot, BrowserCocCocBot, BrowserYahooBot}

var _BrowserNameNameToValueMap = map[string]BrowserName{
	_BrowserNameName[0:14]:         BrowserUnknown,
	_BrowserNameLowerName[0:14]:    BrowserUnknown,
	_BrowserNameName[14:27]:        BrowserChrome,
	_BrowserNameLowerName[14:27]:   BrowserChrome,
	_BrowserNameName[27:36]:        BrowserIE,
	_BrowserNameLowerName[27:36]:   BrowserIE,
	_BrowserNameName[36:49]:        BrowserSafari,
	_BrowserNameLowerName[36:49]:   BrowserSafari,
	_BrowserNameName[49:63]:        BrowserFirefox,
	_BrowserNameLowerName[49:63]:   BrowserFirefox,
	_BrowserNameName[63:77]:        BrowserAndroid,
	_BrowserNameLowerName[63:77]:   BrowserAndroid,
	_BrowserNameName[77:89]:        BrowserOpera,
	_BrowserNameLowerName[77:89]:   BrowserOpera,
	_BrowserNameName[89:106]:       BrowserBlackberry,
	_BrowserNameLowerName[89:106]:  BrowserBlackberry,
	_BrowserNameName[106:122]:      BrowserUCBrowser,
	_BrowserNameLowerName[106:122]: BrowserUCBrowser,
	_BrowserNameName[122:133]:      BrowserSilk,
	_BrowserNameLowerName[122:133]: BrowserSilk,
	_BrowserNameName[133:145]:      BrowserNokia,
	_BrowserNameLowerName[133:145]: BrowserNokia,
	_BrowserNameName[145:160]:      BrowserNetFront,
	_BrowserNameLowerName[145:160]: BrowserNetFront,
	_BrowserNameName[160:169]:      BrowserQQ,
	_BrowserNameLowerName[160:169]: BrowserQQ,
	_BrowserNameName[169:183]:      BrowserMaxthon,
	_BrowserNameLowerName[169:183]: BrowserMaxthon,
	_BrowserNameName[183:203]:      BrowserSogouExplorer,
	_BrowserNameLowerName[183:203]: BrowserSogouExplorer,
	_BrowserNameName[203:217]:      BrowserSpotify,
	_BrowserNameLowerName[203:217]: BrowserSpotify,
	_BrowserNameName[217:232]:      BrowserNintendo,
	_BrowserNameLowerName[217:232]: BrowserNintendo,
	_BrowserNameName[232:246]:      BrowserSamsung,
	_BrowserNameLowerName[232:246]: BrowserSamsung,
	_BrowserNameName[246:259]:      BrowserYandex,
	_BrowserNameLowerName[246:259]: BrowserYandex,
	_BrowserNameName[259:272]:      BrowserCocCoc,
	_BrowserNameLowerName[259:272]: BrowserCocCoc,
	_BrowserNameName[272:282]:      BrowserBot,
	_BrowserNameLowerName[272:282]: BrowserBot,
	_BrowserNameName[282:297]:      BrowserAppleBot,
	_BrowserNameLowerName[282:297]: BrowserAppleBot,
	_BrowserNameName[297:312]:      BrowserBaiduBot,
	_BrowserNameLowerName[297:312]: BrowserBaiduBot,
	_BrowserNameName[312:326]:      BrowserBingBot,
	_BrowserNameLowerName[312:326]: BrowserBingBot,
	_BrowserNameName[326:346]:      BrowserDuckDuckGoBot,
	_BrowserNameLowerName[326:346]: BrowserDuckDuckGoBot,
	_BrowserNameName[346:364]:      BrowserFacebookBot,
	_BrowserNameLowerName[346:364]: BrowserFacebookBot,
	_BrowserNameName[364:380]:      BrowserGoogleBot,
	_BrowserNameLowerName[364:380]: BrowserGoogleBot,
	_BrowserNameName[380:398]:      BrowserLinkedInBot,
	_BrowserNameLowerName[380:398]: BrowserLinkedInBot,
	_BrowserNameName[398:411]:      BrowserMsnBot,
	_BrowserNameLowerName[398:411]: BrowserMsnBot,
	_BrowserNameName[411:428]:      BrowserPingdomBot,
	_BrowserNameLowerName[411:428]: BrowserPingdomBot,
	_BrowserNameName[428:445]:      BrowserTwitterBot,
	_BrowserNameLowerName[428:445]: BrowserTwitterBot,
	_BrowserNameName[445:461]:      BrowserYandexBot,
	_BrowserNameLowerName[445:461]: BrowserYandexBot,
	_BrowserNameName[461:477]:      BrowserCocCocBot,
	_BrowserNameLowerName[461:477]: BrowserCocCocBot,
	_BrowserNameName[477:492]:      BrowserYahooBot,
	_BrowserNameLowerName[477:492]: BrowserYahooBot,
}

var _BrowserNameNames = []string{
	_BrowserNameName[0:14],
	_BrowserNameName[14:27],
	_BrowserNameName[27:36],
	_BrowserNameName[36:49],
	_BrowserNameName[49:63],
	_BrowserNameName[63:77],
	_BrowserNameName[77:89],
	_BrowserNameName[89:106],
	_BrowserNameName[106:122],
	_BrowserNameName[122:133],
	_BrowserNameName[133:145],
	_BrowserNameName[145:160],
	_BrowserNameName[160:169],
	_BrowserNameName[169:183],
	_BrowserNameName[183:203],
	_BrowserNameName[203:217],
	_BrowserNameName[217:232],
	_BrowserNameName[232:246],
	_BrowserNameName[246:259],
	_BrowserNameName[259:272],
	_BrowserNameName[272:282],
	_BrowserNameName[282:297],
	_BrowserNameName[297:312],
	_BrowserNameName[312:326],
	_BrowserNameName[326:346],
	_BrowserNameName[346:364],
	_BrowserNameName[364:380],
	_BrowserNameName[380:398],
	_BrowserNameName[398:411],
	_BrowserNameName[411:428],
	_BrowserNameName[428:445],
	_BrowserNameName[445:461],
	_BrowserNameName[461:477],
	_BrowserNameName[477:492],
}

// BrowserNameString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func BrowserNameString(s string) (BrowserName, error) {
	if val, ok := _BrowserNameNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _BrowserNameNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to BrowserName values", s)
}

// BrowserNameValues returns all values of the enum
func BrowserNameValues() []BrowserName {
	return _BrowserNameValues
}

// BrowserNameStrings returns a slice of all String values of the enum
func BrowserNameStrings() []string {
	strs := make([]string, len(_BrowserNameNames))
	copy(strs, _BrowserNameNames)
	return strs
}

// IsABrowserName returns "true" if the value is listed in the enum definition. "false" otherwise
func (i BrowserName) IsABrowserName() bool {
	for _, v := range _BrowserNameValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for BrowserName
func (i BrowserName) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for BrowserName
func (i *BrowserName) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("BrowserName should be a string, got %s", data)
	}

	var err error
	*i, err = BrowserNameString(s)
	return err
}

// MarshalBSONValue implements the bson.ValueMarshaler interface for BrowserName
func (i BrowserName) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, i.String()), nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for BrowserName
func (i *BrowserName) UnmarshalBSONValue(t bsontype.Type, src []byte) error {
	str, _, _ := bsoncore.ReadString(src)
	var err error
	*i, err = BrowserNameString(str)
	return err
}
