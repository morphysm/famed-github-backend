package devtoolkit_test

import (
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"testing"
)

func TestNewDevToolkit(t *testing.T) {
	tests := []struct {
		name        string
		wantToolkit *devtoolkit.DevToolkit
		wantErr     bool
	}{
		{
			name:        "default_devtoolkit",
			wantToolkit: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := devtoolkit.NewDevToolkit()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevToolkit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
