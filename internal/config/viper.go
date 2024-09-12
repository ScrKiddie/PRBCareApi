package config

import (
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"strings"
)

func NewViper() *viper.Viper {
	config := viper.New()
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Info("Config file not found; using environment variables")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	return config
}
