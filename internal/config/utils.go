package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func bindString(s *string, name string) error {
	envValue, ok := loadStringFromEnvironment(name)
	if ok != nil {
		return ok
	}

	*s = envValue
	return nil
}

func loadStringFromEnvironment(name string) (string, error) {
	viper.BindEnv(name)
	envValue := viper.GetString(name)
	if envValue == "" {
		return "", fmt.Errorf("no %s environment variable found", name)
	}
	return envValue, nil
}
