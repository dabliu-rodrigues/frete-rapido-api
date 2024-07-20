package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jsGolden/frete-rapido-api/config"
	"github.com/jsGolden/frete-rapido-api/router"
)

func main() {
	err := config.SetupEnvFile()
	if err != nil {
		log.Fatalf("Error while loading .env file: %s", err)
	}

	serverConfig := config.ServerConfig()

	quoteMongoService, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	router := router.SetupRouter(quoteMongoService)

	fmt.Printf("Server listening at: %s 🚀", serverConfig)
	if err := http.ListenAndServe(serverConfig, router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
