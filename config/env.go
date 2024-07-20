package config

import (
	"github.com/spf13/viper"
)

func SetupEnvFile() error {
	viper.SetConfigFile(".env")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("MONGO_URL", "mongodb://localhost:27017/")
	viper.SetDefault("MONGO_DATABASE", "frete-rapido")
	viper.SetDefault("FRETE_RAPIDO_API_URL", "https://sp.freterapido.com/api/v3")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return err
		}
	}

	return nil
}
