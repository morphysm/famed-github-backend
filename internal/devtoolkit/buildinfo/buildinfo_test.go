package buildinfo

import (
	"github.com/Masterminds/semver/v3"
	"reflect"
	"testing"
	"time"
)

func TestBuildInfo_UserAgent(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := BuildInfo{
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
		want    *BuildInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuildInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuildInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuildInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
