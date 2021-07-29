package uasurfer

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_Version(t *testing.T) {
	type test struct {
		Version     Version
		OSName      OSName
		BrowserName BrowserName
		DeviceType  DeviceType
		Platform    Platform
	}

	v := test{
		Version: Version{
			1,
			2,
			3,
		},
		OSName:      1,
		BrowserName: 1,
		DeviceType:  1,
		Platform:    1,
	}

	// Marshal the struct and test if all enums are stringified.
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"Version":"1.2.3","OSName":"OSWindowsPhone","BrowserName":"BrowserChrome","DeviceType":"DeviceComputer","Platform":"PlatformWindows"}`
	if string(b) != want {
		t.Fatalf("want = %s, got = %s", want, string(b))
	}

	// Create a new struct and check if the Marshaled json is the same after Unmarshal.
	var s2 test
	if err := json.Unmarshal(b, &s2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s2, v) {
		t.Fatal("not equal")
	}
}
