package server_test

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/health"
	"github.com/morphysm/famed-github-backend/internal/server"
	"reflect"
	"testing"
)

func TestFamedAdminRoutes(t *testing.T) {
	type args struct {
		g             *echo.Group
		famedHandler  famed.HTTPHandler
		githubHandler github.HTTPHandler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.FamedAdminRoutes(tt.args.g, tt.args.famedHandler, tt.args.githubHandler)
		})
	}
}

func TestFamedRoutes(t *testing.T) {
	type args struct {
		g       *echo.Group
		handler famed.HTTPHandler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.FamedRoutes(tt.args.g, tt.args.handler)
		})
	}
}

func TestHealthRoutes(t *testing.T) {
	type args struct {
		g       *echo.Group
		handler health.HTTPHandler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.HealthRoutes(tt.args.g, tt.args.handler)
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		devToolKit *devtoolkit.DevToolkit
	}
	tests := []struct {
		name    string
		args    args
		want    *server.Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.NewServer(tt.args.devToolKit)
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
