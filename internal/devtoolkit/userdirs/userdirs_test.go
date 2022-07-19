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
		{
			name: "no_path",
			args: args{
				programName: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserDirs(tt.args.programName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.ConfigHome == "" {
				t.Errorf("NewUserDirs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
