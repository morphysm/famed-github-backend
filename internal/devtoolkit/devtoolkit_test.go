package devtoolkit

import (
	"testing"
)

func TestNewDevToolkit(t *testing.T) {
	tests := []struct {
		name        string
		wantToolkit *DevToolkit
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
			_, err := NewDevToolkit()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevToolkit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
