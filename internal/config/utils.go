package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func bindBool(s *bool, name string) error {
	envValue, ok := loadBoolFromEnvironment(name)
	if ok != nil {
		return ok
	}

	*s = envValue

	return nil
}

func loadBoolFromEnvironment(name string) (bool, error) {
	if err := viper.BindEnv(name); err != nil {
		return false, err
	}

	envValue := viper.GetBool(name)
	return envValue, nil
}

func bindString(s *string, name string) error {
	envValue, ok := loadStringFromEnvironment(name)
	if ok != nil {
		return ok
	}

	*s = envValue

	return nil
}

func loadStringFromEnvironment(name string) (string, error) {
	if err := viper.BindEnv(name); err != nil {
		return "", err
	}

	envValue := viper.GetString(name)
	if envValue == "" {
		return "", fmt.Errorf("no %s environment variable found", name)
	}

	return envValue, nil
}

func bindInt64(i *int64, name string) error {
	envValue, ok := loadInt64FromEnvironment(name)
	if ok != nil {
		return ok
	}

	*i = envValue

	return nil
}

func loadInt64FromEnvironment(name string) (int64, error) {
	if err := viper.BindEnv(name); err != nil {
		return 0, err
	}

	envValue := viper.GetInt64(name)
	if envValue == 0 {
		return 0, fmt.Errorf("no %s environment variable found", name)
	}

	return envValue, nil
}
