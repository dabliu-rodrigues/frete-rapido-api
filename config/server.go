package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port string
}

func ServerConfig() string {
	server := fmt.Sprintf("%s:%s", viper.GetString("HOST"), viper.GetString("PORT"))
	return server
}
