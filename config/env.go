package config

import (
	"github.com/spf13/viper"
)

func SetupEnvFile() error {
	viper.SetConfigFile(".env")
	viper.SetDefault("PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return err
		}
	}

	return nil
}
