package server

import (
	"github.com/labstack/echo/v4"
	"github.com/morphysm/famed-github-backend/internal/config"
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/famed"
	"github.com/morphysm/famed-github-backend/internal/github"
	"github.com/morphysm/famed-github-backend/internal/health"
	"github.com/newrelic/go-agent/v3/newrelic"
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
			FamedAdminRoutes(tt.args.g, tt.args.famedHandler, tt.args.githubHandler)
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
			FamedRoutes(tt.args.g, tt.args.handler)
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
			HealthRoutes(tt.args.g, tt.args.handler)
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
		want    *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServer(tt.args.devToolKit)
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
		echo       *echo.Echo
		devToolKit *devtoolkit.DevToolkit
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
			s := &Server{
				echo:       tt.fields.echo,
				devToolKit: tt.fields.devToolKit,
			}
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_configureNewRelic(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *newrelic.Application
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := configureNewRelic(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("configureNewRelic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configureNewRelic() got = %v, want %v", got, tt.want)
			}
		})
	}
}
