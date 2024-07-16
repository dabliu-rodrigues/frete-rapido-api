package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/jsGolden/frete-rapido-api/config"
)

func main() {
	err := config.SetupEnvFile()
	if err != nil {
		log.Fatalf("Error while loading .env file: %s", err)
	}

	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	port := viper.GetString("PORT")

	fmt.Printf("Server listening at port :%s 🚀", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
