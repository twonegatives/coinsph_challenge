package config

import "github.com/spf13/viper"

type configDefaults struct {
	Listen string
	AppEnv string
}

func getDefaults() *configDefaults {
	return &configDefaults{
		Listen: ":80",
		AppEnv: "dev",
	}
}

func NewConfig() *viper.Viper {
	defaults := getDefaults()
	cfg := viper.New()

	cfg.SetDefault("LISTEN", defaults.Listen)
	cfg.SetDefault("APP_ENV", defaults.AppEnv)
	cfg.SetDefault("SHUTDOWN_TIMEPUT", "2s")
	cfg.AutomaticEnv()

	return cfg
}
