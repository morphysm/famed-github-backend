package config

import (
	"github.com/knadh/koanf"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"testing"
)

func Test_loadConfigFile(t *testing.T) {
	type args struct {
		koanf    *koanf.Koanf
		filePath string
		parser   koanf.Parser
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadConfigFile(tt.args.koanf, tt.args.filePath, tt.args.parser); (err != nil) != tt.wantErr {
				t.Errorf("loadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_loadDefaultValues(t *testing.T) {
	type args struct {
		koanf *koanf.Koanf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadDefaultValues(tt.args.koanf); (err != nil) != tt.wantErr {
				t.Errorf("loadDefaultValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_loadEnvVars(t *testing.T) {
	type args struct {
		koanf *koanf.Koanf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadEnvVars(tt.args.koanf); (err != nil) != tt.wantErr {
				t.Errorf("loadEnvVars() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("verifyConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyLabel(t *testing.T) {
	type args struct {
		cfg   Config
		label string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyLabel(tt.args.cfg, tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("verifyLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyReward(t *testing.T) {
	type args struct {
		cfg  Config
		cvss model.IssueSeverity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyReward(tt.args.cfg, tt.args.cvss); (err != nil) != tt.wantErr {
				t.Errorf("verifyReward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
