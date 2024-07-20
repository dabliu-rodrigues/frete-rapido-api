package config

import (
	"github.com/jsGolden/frete-rapido-api/services"
	"github.com/spf13/viper"
)

func SetupDatabase() (*services.MongoService, error) {
	mongoService := services.NewMongoService(viper.GetString("MONGO_URL"), viper.GetString("MONGO_DATABASE"))

	_, err := mongoService.GetConnection()
	if err != nil {
		return nil, err
	}
	return mongoService, nil
}
