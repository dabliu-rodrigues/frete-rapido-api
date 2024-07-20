package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jsGolden/frete-rapido-api/handlers/quotes"
	"github.com/jsGolden/frete-rapido-api/middlewares"
	"github.com/jsGolden/frete-rapido-api/services"
	"github.com/spf13/viper"
)

func SetupRouter(q *services.MongoService) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.Cors())

	freteRapidoService := services.NewFreteRapidoService(viper.GetString("FRETE_RAPIDO_API_URL"))
	quoteHandler := quotes.NewQuoteHandler(q, freteRapidoService)

	router.Post("/quote", quoteHandler.CreateQuote)
	router.Get("/metrics", quotes.GetQuoteMetrics)
	return router
}
