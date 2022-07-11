package devtoolkit

import (
	"reflect"
	"testing"
)

func TestNewDevToolkit(t *testing.T) {
	tests := []struct {
		name        string
		wantToolkit *DevToolkit
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToolkit, err := NewDevToolkit()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevToolkit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotToolkit, tt.wantToolkit) {
				t.Errorf("NewDevToolkit() gotToolkit = %v, want %v", gotToolkit, tt.wantToolkit)
			}
		})
	}
}
