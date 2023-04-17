package uasurfer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Version_findVersionNumber(t *testing.T) {
	type args struct {
		ua            string
		versionPrefix string
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantVersion Version
	}{
		{
			name: "major only",
			args: args{
				ua:            "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30 Safari/537.36",
				versionPrefix: "Chrome/",
			},
			want: true,
			wantVersion: Version{
				Major: 30,
			},
		},
		{
			name: "major and minor",
			args: args{
				ua:            "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.1 Safari/537.36",
				versionPrefix: "Chrome/",
			},
			want: true,
			wantVersion: Version{
				Major: 30,
				Minor: 1,
			},
		},
		{
			name: "major, minor and patch",
			args: args{
				ua:            "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.1.2 Safari/537.36",
				versionPrefix: "Chrome/",
			},
			want: true,
			wantVersion: Version{
				Major: 30,
				Minor: 1,
				Patch: 2,
			},
		},
		{
			name: "get label and trim trailing special characters",
			args: args{
				ua:            "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.1.2.3; Safari/537.36",
				versionPrefix: "Chrome/",
			},
			want: true,
			wantVersion: Version{
				Major: 30,
				Minor: 1,
				Patch: 2,
				Extra: "3",
			},
		},
		{
			name: "get complex label",
			args: args{
				ua:            "Mozilla/5.0 (Nintendo WiiU) AppleWebKit/536.30 (KHTML, like Gecko) NX/3.0.4.2.12 NintendoBrowser/4.3.1.11264.US",
				versionPrefix: "NintendoBrowser/",
			},
			want: true,
			wantVersion: Version{
				Major: 4,
				Minor: 3,
				Patch: 1,
				Extra: "11264.US",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{}
			if got := v.findVersionNumber(tt.args.ua, tt.args.versionPrefix); got != tt.want {
				t.Fatalf("Version.findVersionNumber() = %v, want %v", got, tt.want)
			}
			if diff := cmp.Diff(v, tt.wantVersion); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
