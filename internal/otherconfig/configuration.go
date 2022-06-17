package otherconfig

import (
	"github.com/awnumar/memguard"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/rotisserie/eris"
	"os"
	"strings"
)

// delimiter ..
const delimiter = "."

// envPrefix ..
const envPrefix = "FAMED_"

type Config struct {
	App struct {
		Host string `koanf:"host"`
		Port string `koanf:"port"`
	} `koanf:"app"`

	NewRelic struct {
		Name    string `koanf:"name"`
		Key     string `koanf:"key"`
		Enabled bool   `koanf:"enabled"`
	} `koanf:"newrelic"`

	Github struct {
		Host          string            `koanf:"host"`
		KeyEnclave    *memguard.Enclave `koanf:"keyenclave"`
		WebhookSecret string            `koanf:"webhooksecret"`
		AppID         int64             `koanf:"appid"`
		BotLogin      string            `koanf:"botlogin"`
	} `koanf:"github"`

	Famed struct {
		Labels          map[string]model.Label          `koanf:"labels"`
		Rewards         map[model.IssueSeverity]float64 `koanf:"rewards"`
		Currency        string                          `koanf:"currency"`
		DaysToFix       int                             `koanf:"daystofix"`
		UpdateFrequency int                             `koanf:"updatefrequency"`
	} `koanf:"famed"`

	// TODO this should probably not be in memory
	RedTeamLogins map[string]string `koanf:"redteamslogins"`

	Admin struct {
		Username string `koanf:"username"`
		Password string `koanf:"password"`
	} `koanf:"admin"`
}

func NewConfig(filePath string) (*Config, error) {
	k := koanf.New(delimiter)

	// TODO currency.host don't exist. what to do ?

	defaultLabel := map[string]model.Label{}

	defaultLabel["famed"] = model.Label{
		Name:        "famed",
		Color:       "566FDB",
		Description: "Famed - Tracked by Famed",
	}

	defaultLabel["info"] = model.Label{
		Name:        "info",
		Color:       "566FDB",
		Description: "Famed - Common Vulnerability Scoring System (CVSS) - None",
	}

	defaultLabel["low"] = model.Label{
		Name:        "low",
		Color:       "566FDB",
		Description: "Famed - Common Vulnerability Scoring System (CVSS) - Low",
	}

	defaultLabel["medium"] = model.Label{
		Name:        "medium",
		Color:       "566FDB",
		Description: "Famed - Common Vulnerability Scoring System (CVSS) - Medium",
	}

	defaultLabel["high"] = model.Label{
		Name:        "high",
		Color:       "566FDB",
		Description: "Famed - Common Vulnerability Scoring System (CVSS) - High",
	}

	defaultLabel["critical"] = model.Label{
		Name:        "critical",
		Color:       "566FDB",
		Description: "Famed - Common Vulnerability Scoring System (CVSS) - Critical",
	}

	defaultIssueSeverity := map[model.IssueSeverity]float64{}

	defaultIssueSeverity[model.Info] = 0
	defaultIssueSeverity[model.Low] = 1000
	defaultIssueSeverity[model.Medium] = 5000
	defaultIssueSeverity[model.High] = 10000
	defaultIssueSeverity[model.Critical] = 25000

	defaultRedTeams := map[string]string{}
	defaultRedTeams["Jonny Rhea"] = "jrhea"
	defaultRedTeams["Alexander Sadovskyi"] = "AlexSSD7"
	defaultRedTeams["Martin Holst Swende"] = "holiman"
	defaultRedTeams["Tintin"] = "tintinweb"
	defaultRedTeams["Antoine Toulme"] = "atoulme"
	defaultRedTeams["Stefan Kobrc"] = "tintinweb"
	defaultRedTeams["Quan"] = "cryptosubtlety"
	defaultRedTeams["WINE Academic Workshop"] = ""
	defaultRedTeams["Proto"] = "protolambda"
	defaultRedTeams["Taurus"] = ""
	defaultRedTeams["Saulius Grigaitis (+team)."] = "sifraitech"
	defaultRedTeams["Antonio Sanso"] = "asanso"
	defaultRedTeams["Guido Vranken"] = "guidovranken"
	defaultRedTeams["Guido Vranken"] = "guidovranken"
	defaultRedTeams["Jacek"] = "arnetheduck"
	defaultRedTeams["Onur Kılıç"] = "kilic"
	defaultRedTeams["Jim McDonald"] = "mcdee"
	defaultRedTeams["Nishant (Prysm)"] = "nisdas"

	// Load defaults values.
	err := k.Load(confmap.Provider(map[string]interface{}{
		"app.host":              "127.0.0.1",
		"app.port":              "8080",
		"github.host":           "https://api.github.com",
		"currency.host":         "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1",
		"famed.labels":          defaultLabel,
		"famed.rewards":         defaultIssueSeverity,
		"famed.currency":        "POINTS",
		"famed.daystofix":       90,
		"famed.updatefrequency": 120,
		"famed.redteamslogins":  defaultRedTeams,
	}, "."), nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to load default values")
	}

	if yamlFile, err := os.Stat(filePath); err == nil {
		// Load YAML config from filePath
		if err := k.Load(file.Provider(yamlFile.Name()), yaml.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load yaml config")
		}
	}

	if envFile, err := os.Stat(".env"); err == nil {
		// Load from ./.env file
		if err := k.Load(file.Provider(envFile.Name()), dotenv.Parser()); err != nil {
			return nil, eris.Wrap(err, "failed to load dotenv config")
		}
	}

	// Load from environment variables
	err = k.Load(env.Provider(envPrefix, delimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to load environment variables")
	}

	var config *Config

	err = k.Unmarshal("", &config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal config")
	}

	return config, nil
}
