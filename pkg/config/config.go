// Package config provides functionality to alter application
// behaviour based on a set of passed environment variables
package config

import "github.com/spf13/viper"

type configDefaults struct {
	Listen string
	AppEnv string
	DB     string
}

func getDefaults() *configDefaults {
	return &configDefaults{
		Listen: ":80",
		AppEnv: "dev",
		DB:     "postgres://localhost/coinsph?sslmode=disable",
	}
}

func NewConfig() *viper.Viper {
	defaults := getDefaults()
	cfg := viper.New()

	cfg.SetDefault("LISTEN", defaults.Listen)
	cfg.SetDefault("APP_ENV", defaults.AppEnv)
	cfg.SetDefault("DB", defaults.DB)
	cfg.SetDefault("SHUTDOWN_TIMEOUT", "2s")
	cfg.AutomaticEnv()

	return cfg
}
