package config

import (
	"github.com/morphysm/famed-github-backend/internal/repositories/github/model"
	"github.com/phuslu/log"
)

var defaultConfig = map[string]interface{}{
	"app.host":      "127.0.0.1",
	"app.port":      "8080",
	"app.loglevel":  log.ErrorLevel,
	"github.host":   "https://api.github.com",
	"currency.host": "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1", // TODO: NOT EXIST IN ORIGINAL SOURCE ?!
	"famed.labels": map[string]model.Label{
		"famed": {
			Name:        "famed",
			Color:       "566FDB",
			Description: "Famed - Tracked by Famed",
		},
		"info": {
			Name:        "info",
			Color:       "566FDB",
			Description: "Fam²ed - Common Vulnerability Scoring System (CVSS) - None",
		},
		"low": {
			Name:        "low",
			Color:       "566FDB",
			Description: "Famed - Common Vulnerability Scoring System (CVSS) - Low",
		},
		"medium": {
			Name:        "medium",
			Color:       "566FDB",
			Description: "Famed - Common Vulnerability Scoring System (CVSS) - Medium",
		},
		"high": {
			Name:        "high",
			Color:       "566FDB",
			Description: "Famed - Common Vulnerability Scoring System (CVSS) - High",
		},
		"critical": {
			Name:        "critical",
			Color:       "566FDB",
			Description: "Famed - Common Vulnerability Scoring System (CVSS) - Critical",
		},
	},
	"famed.rewards": map[model.IssueSeverity]float64{
		model.Info:     0,
		model.Low:      1000,
		model.Medium:   5000,
		model.High:     10000,
		model.Critical: 25000,
	},
	"famed.currency":        "POINTS",
	"famed.daystofix":       90,
	"famed.updatefrequency": 120,
	"famed.redteamlogins": map[string]string{
		"Jonny Rhea":                 "jrhea",
		"Alexander Sadovskyi":        "AlexSSD7",
		"Martin Holst Swende":        "holiman",
		"Tintin":                     "tintinweb",
		"Antoine Toulme":             "atoulme",
		"Stefan Kobrc":               "",
		"Quan":                       "cryptosubtlety",
		"WINE Academic Workshop":     "",
		"Proto":                      "protolambda",
		"Taurus":                     "",
		"Saulius Grigaitis (+team).": "sifraitech",
		"Antonio Sanso":              "asanso",
		"Guido Vranken":              "guidovranken",
		"Jacek":                      "arnetheduck",
		"Onur Kılıç":                 "kilic",
		"Jim McDonald":               "mcdee",
		"Nishant (Prysm)":            "nisdas",
	},
}
