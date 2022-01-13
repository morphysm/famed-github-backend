package config

import (
	"fmt"

	"github.com/spf13/cast"
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

func bindInt64(i *int64, name string) error {
	envValue, ok := loadInt64FromEnvironment(name)
	if ok != nil {
		return ok
	}

	*i = envValue

	return nil
}

func loadInt64FromEnvironment(name string) (int64, error) {
	viper.BindEnv(name)
	envValue := viper.GetInt64(name)
	if envValue == 0 {
		return 0, fmt.Errorf("no %s environment variable found", name)
	}

	return envValue, nil
}

func bindInt64Slice(s *[]int64, name string) error {
	envValue, ok := loadIntSliceFromEnvironment(name)
	if ok != nil {
		return ok
	}

	var envValue64 []int64
	for _, v := range envValue {
		envValue64 = append(envValue64, int64(v))
	}

	*s = envValue64

	return nil
}

func loadIntSliceFromEnvironment(name string) ([]int, error) {
	viper.BindEnv(name)
	envValue := viper.GetStringSlice(name)
	if len(envValue) == 0 {
		return nil, fmt.Errorf("no %s environment variable found", name)
	}
	// Needs this workaround, viper.GetIntSlice is buggy.
	intSlice := cast.ToIntSlice(envValue)
	if len(intSlice) == 0 {
		return nil, fmt.Errorf("no %s environment variable found", name)
	}

	return intSlice, nil
}
