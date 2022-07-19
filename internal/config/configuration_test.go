package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	f, _ := os.Create("config_test.json")
	defer f.Close()

	f.WriteString("{\n  \"github\": {\n    \"key\": \"-----BEGIN RSA PRIVATE KEY-----\\nMIIEpgIBAAKCAQEA7KhFsYbhfcaq6iqFTJYYaNXdE32PMvN96QOMUHLNwEmWr+C3\\nMCfmMCrV0wSUKwh/YnjIGiQWNQjWuZdG9Ic3y/cOOrpsORLYEjw43qLMvojaWx43\\nZ9mQH3lpwgqhbk/OoS7l8tNM7q+7Ua/haCYQYmvOUCv+CyPaVfeYp28h75Hw/CUJ\\nbqTlOezWSPWGOdYKN97pQcVTDRKdUmSeNEEGcBkO6pZ4Ah7sNKYBMANJRKmAgdPi\\nYpyv/qYNsHj/PXlnrsZ2/ayk/jmlg7SA/NABi4SC+D9i6wBOloOeOSrSw+q9ObsX\\nwfFFWA3K8cw+VKdwf8Q14jtfeD2GHj87K1QHhwIDAQABAoIBAQCmBnNGVSLykyKq\\nvwPfM9mSCp9bIhYJH6twgk24zqGrybSOVK8PeJ5Toml57ddozUBYu/Vd6X0u3bGO\\naCOePxKU5BC2gLyV2bN+L4OSJVJQRUAy9mLWV1p1yj64o66W7iQ/DeDCVxy8wso+\\nR45x+2o5Mfp+Yi6KcC+nadlNdXiwUTAICWPOMZgj9AwJGA7lniopEU7HeKYUqyEW\\nEM6vZHXOudIm/Yboaw0PJN1Z7Y/bOIFY2oo3ow6t7ahcrsohs4yZhvr05ZXPXoSC\\nKMMXDrvLBTTelL0zJN4f0UhLcvsYPQJJmneMG96nmr33OQaWyboEXHdrip4qsjj7\\nNdM4dHJpAoGBAPiO5Np9x2VC6nQ22Gb0tVfvrnA3Fm2EFyjcxMRmkI2pYJVBnxTR\\nwZCvJ+csHVhNN+39SYL9X5PKklahH+0nNsV45VGPeB0lCKWBWQsOrNWH8hTI/r5d\\nbqyjS1BKY8we2RXP0ZPVrX7SCrmFQfIZJdCc3g/Ve5gaE9y9m7A4XdLlAoGBAPO+\\nKZ38YKyvG9FaH8t6noGzK5oHHgZ5PvOb+Mo01fRXHh6pyWEvS7zZWkkVK66DgTRB\\nSnObOA17UVWj6waUqrWDqUYxe64iXT/nYPi6qYXntxzukjtfxY2AJM9rCAvJ1XpG\\nreGVqc4LEA+DrOG2SNqzXL+PPwU3JrIYNhtb9C37AoGBANMlyGGXgeCStMqeoLzt\\nWnPmR0BKe8Hy+R2cVYcmPdwpq8N/aF1uRsnbEcG+5vrRNhb1GRKunRfWePQgkheL\\nPWsJZX0grH/Nqwe11ueewtHuV4ayrD0Y7+C2I0+Esjx/ZBi0Xyv/1A+s7LFm83tv\\nQ4FxEO9QglrWpFLbu7s6VvHFAoGBAJJ5vNjMSex8buMoneLSFV8sJQ+zJ0AMrOAI\\n40Hg7pKfp+IVdoeIvKMIm1E//7goHwUgF3XR2aWAbihhEWQrA0uBi8A7DHBhBljY\\n21WeFzH5RfmFBSvZKgcW8wgS8grjh/6rauMd5aWE0GoCX2pk+PM0xo/3rY+czQxJ\\nsHpQkDTxAoGBAN625+aJjFR8A9jikwHPCwni/bKXutSv3HZkSb/00F4A1tkicozQ\\ncsaq42Us3Ri5TUJwNgg4balHhWWfC9Y3hZ7dukw8XD5MDNByM+YgzcpAPGZaKoLK\\nfcm/vHi5UZMyvMksindgYqcO8RczL+7t73qrMHGkDZ/AyRB7C9UcJTNJ\\n-----END RSA PRIVATE KEY-----\\n\",\n    \"webhooksecret\": \"foobar\",\n    \"appid\": 1234,\n    \"botlogin\": \"foobar\"\n  },\n  \"admin\": {\n    \"username\": \"foobar\",\n    \"password\": \"foobar\"\n  }\n}")

	type args struct {
		filePath string
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "without_config",
			args: args{
				filePath: "noconfig.json",
			},
			wantConfig: nil,
			wantErr:    true,
		},
		{
			name: "normal",
			args: args{
				filePath: "config_test.json",
			},
			wantConfig: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := NewConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotConfig.Admin.Password != "foobar" {
				t.Errorf("NewConfig() gotConfig = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}

func Test_loadConfigFile(t *testing.T) {
	k := koanf.New(delimiter)

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
		{
			name: "without_config",
			args: args{
				koanf:    k,
				filePath: "noconfig.json",
				parser:   json.Parser(),
			},
			wantErr: true,
		},
		{
			name: "with_config",
			args: args{
				koanf:    k,
				filePath: "config_test.json",
				parser:   json.Parser(),
			},
			wantErr: false,
		},
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
	k := koanf.New(delimiter)

	type args struct {
		koanf *koanf.Koanf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "load_default",
			args: args{
				koanf: k,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadDefaultValues(tt.args.koanf); (err != nil) != tt.wantErr {
				t.Errorf("loadDefaultValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Todo: find a way to test env var ?
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

func Test_verifyLabel(t *testing.T) {
	cfg := Config{}

	k := koanf.New(delimiter)
	loadDefaultValues(k)

	k.Unmarshal("", &cfg)

	type args struct {
		cfg   Config
		label string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exist",
			args: args{
				cfg:   cfg,
				label: string(model.Info),
			},
			wantErr: false,
		},
		{
			name: "not_exist",
			args: args{
				cfg:   cfg,
				label: "foo",
			},
			wantErr: true,
		},
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
	cfg := Config{}

	cfg2 := Config{}

	k := koanf.New(delimiter)
	loadDefaultValues(k)

	k.Unmarshal("", &cfg2)

	type args struct {
		cfg  Config
		cvss model.IssueSeverity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				cfg:  cfg,
				cvss: "",
			},
			wantErr: true,
		},
		{
			name: "default",
			args: args{
				cfg:  cfg2,
				cvss: model.Info,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyReward(tt.args.cfg, tt.args.cvss); (err != nil) != tt.wantErr {
				t.Errorf("verifyReward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyConfig(t *testing.T) {
	cfg := Config{}

	cfg2 := Config{}

	k := koanf.New(delimiter)
	loadDefaultValues(k)

	cfg2.Github.Key = "fookey"
	cfg2.Github.WebhookSecret = "webhook"
	cfg2.Github.AppID = 1337
	cfg2.Github.BotLogin = "bot"
	cfg2.Admin.Username = "foousername"
	cfg2.Admin.Password = "foopassword"

	k.Unmarshal("", &cfg2)

	type args struct {
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty_struct",
			args: args{
				cfg: cfg,
			},
			wantErr: true,
		},
		{
			name: "full_struct",
			args: args{
				cfg: cfg2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("verifyConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}