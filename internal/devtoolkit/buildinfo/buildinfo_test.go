package buildinfo_test

import (
	"github.com/Masterminds/semver/v3"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit/buildinfo"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestBuildInfo_UserAgent(t *testing.T) {
	version, _ := semver.NewVersion("0.0.0")

	type fields struct {
		Version         *semver.Version
		Date            time.Time
		Revision        string
		Target          string
		CompilerVersion string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Version:         version,
				Date:            time.Time{},
				Revision:        "rev",
				Target:          "tar",
				CompilerVersion: "compvers",
			},
			want: "famed/0.0.0 (" + runtime.Version() + "; tar)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := buildinfo.BuildInfo{
				Version:         tt.fields.Version,
				Date:            tt.fields.Date,
				Revision:        tt.fields.Revision,
				Target:          tt.fields.Target,
				CompilerVersion: tt.fields.CompilerVersion,
			}
			if got := i.UserAgent(); got != tt.want {
				t.Errorf("UserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBuildInfo(t *testing.T) {
	tests := []struct {
		name    string
		want    *buildinfo.BuildInfo
		wantErr bool
	}{
		{
			name:    "default_buildinfo",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildinfo.NewBuildInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuildInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.CompilerVersion, runtime.Version()) {
				t.Errorf("NewBuildInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
