package userdirs

import (
	"testing"
)

func TestNewUserDirs(t *testing.T) {
	type args struct {
		programName string
	}
	tests := []struct {
		name    string
		args    args
		want    *UserDirs
		wantErr bool
	}{
		{
			name: "default_path",
			args: args{
				programName: "testapp",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserDirs(tt.args.programName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ConfigHome == "" {
				t.Errorf("NewUserDirs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDirs_makePaths(t *testing.T) {
	type fields struct {
		ConfigHome string
		CacheHome  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserDirs{
				ConfigHome: tt.fields.ConfigHome,
				CacheHome:  tt.fields.CacheHome,
			}
			if err := u.makePaths(); (err != nil) != tt.wantErr {
				t.Errorf("makePaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
