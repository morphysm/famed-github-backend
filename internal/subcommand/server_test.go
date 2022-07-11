package subcommand

import (
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/server"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	type args struct {
		devtoolkit *devtoolkit.DevToolkit
	}
	tests := []struct {
		name    string
		args    args
		want    *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServer(tt.args.devtoolkit)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	type fields struct {
		DevToolkit *devtoolkit.DevToolkit
		Server     *server.Server
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
			a := &Server{
				DevToolkit: tt.fields.DevToolkit,
				Server:     tt.fields.Server,
			}
			if err := a.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
