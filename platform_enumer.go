// Code generated by "enumer -type=Platform -json -bson"; DO NOT EDIT.

package uasurfer

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strings"
)

const _PlatformName = "PlatformUnknownPlatformWindowsPlatformMacPlatformLinuxPlatformiPadPlatformiPhonePlatformiPodPlatformBlackberryPlatformWindowsPhonePlatformPlaystationPlatformXboxPlatformNintendoPlatformBot"

var _PlatformIndex = [...]uint8{0, 15, 30, 41, 54, 66, 80, 92, 110, 130, 149, 161, 177, 188}

const _PlatformLowerName = "platformunknownplatformwindowsplatformmacplatformlinuxplatformipadplatformiphoneplatformipodplatformblackberryplatformwindowsphoneplatformplaystationplatformxboxplatformnintendoplatformbot"

func (i Platform) String() string {
	if i < 0 || i >= Platform(len(_PlatformIndex)-1) {
		return fmt.Sprintf("Platform(%d)", i)
	}
	return _PlatformName[_PlatformIndex[i]:_PlatformIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PlatformNoOp() {
	var x [1]struct{}
	_ = x[PlatformUnknown-(0)]
	_ = x[PlatformWindows-(1)]
	_ = x[PlatformMac-(2)]
	_ = x[PlatformLinux-(3)]
	_ = x[PlatformiPad-(4)]
	_ = x[PlatformiPhone-(5)]
	_ = x[PlatformiPod-(6)]
	_ = x[PlatformBlackberry-(7)]
	_ = x[PlatformWindowsPhone-(8)]
	_ = x[PlatformPlaystation-(9)]
	_ = x[PlatformXbox-(10)]
	_ = x[PlatformNintendo-(11)]
	_ = x[PlatformBot-(12)]
}

var _PlatformValues = []Platform{PlatformUnknown, PlatformWindows, PlatformMac, PlatformLinux, PlatformiPad, PlatformiPhone, PlatformiPod, PlatformBlackberry, PlatformWindowsPhone, PlatformPlaystation, PlatformXbox, PlatformNintendo, PlatformBot}

var _PlatformNameToValueMap = map[string]Platform{
	_PlatformName[0:15]:         PlatformUnknown,
	_PlatformLowerName[0:15]:    PlatformUnknown,
	_PlatformName[15:30]:        PlatformWindows,
	_PlatformLowerName[15:30]:   PlatformWindows,
	_PlatformName[30:41]:        PlatformMac,
	_PlatformLowerName[30:41]:   PlatformMac,
	_PlatformName[41:54]:        PlatformLinux,
	_PlatformLowerName[41:54]:   PlatformLinux,
	_PlatformName[54:66]:        PlatformiPad,
	_PlatformLowerName[54:66]:   PlatformiPad,
	_PlatformName[66:80]:        PlatformiPhone,
	_PlatformLowerName[66:80]:   PlatformiPhone,
	_PlatformName[80:92]:        PlatformiPod,
	_PlatformLowerName[80:92]:   PlatformiPod,
	_PlatformName[92:110]:       PlatformBlackberry,
	_PlatformLowerName[92:110]:  PlatformBlackberry,
	_PlatformName[110:130]:      PlatformWindowsPhone,
	_PlatformLowerName[110:130]: PlatformWindowsPhone,
	_PlatformName[130:149]:      PlatformPlaystation,
	_PlatformLowerName[130:149]: PlatformPlaystation,
	_PlatformName[149:161]:      PlatformXbox,
	_PlatformLowerName[149:161]: PlatformXbox,
	_PlatformName[161:177]:      PlatformNintendo,
	_PlatformLowerName[161:177]: PlatformNintendo,
	_PlatformName[177:188]:      PlatformBot,
	_PlatformLowerName[177:188]: PlatformBot,
}

var _PlatformNames = []string{
	_PlatformName[0:15],
	_PlatformName[15:30],
	_PlatformName[30:41],
	_PlatformName[41:54],
	_PlatformName[54:66],
	_PlatformName[66:80],
	_PlatformName[80:92],
	_PlatformName[92:110],
	_PlatformName[110:130],
	_PlatformName[130:149],
	_PlatformName[149:161],
	_PlatformName[161:177],
	_PlatformName[177:188],
}

// PlatformString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PlatformString(s string) (Platform, error) {
	if val, ok := _PlatformNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PlatformNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Platform values", s)
}

// PlatformValues returns all values of the enum
func PlatformValues() []Platform {
	return _PlatformValues
}

// PlatformStrings returns a slice of all String values of the enum
func PlatformStrings() []string {
	strs := make([]string, len(_PlatformNames))
	copy(strs, _PlatformNames)
	return strs
}

// IsAPlatform returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Platform) IsAPlatform() bool {
	for _, v := range _PlatformValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Platform
func (i Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Platform
func (i *Platform) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Platform should be a string, got %s", data)
	}

	var err error
	*i, err = PlatformString(s)
	return err
}

// MarshalBSONValue implements the bson.ValueMarshaler interface for Platform
func (i Platform) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, i.String()), nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for Platform
func (i *Platform) UnmarshalBSONValue(t bsontype.Type, src []byte) error {
	str, _, _ := bsoncore.ReadString(src)
	var err error
	*i, err = PlatformString(str)
	return err
}
