package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jsGolden/frete-rapido-api/handlers/quotes"
	"github.com/jsGolden/frete-rapido-api/middlewares"
	"github.com/jsGolden/frete-rapido-api/repositories"
	"github.com/jsGolden/frete-rapido-api/services"
	freterapido "github.com/jsGolden/frete-rapido-api/services/frete-rapido"
	"github.com/spf13/viper"
)

func SetupRouter(m *services.MongoService) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.Cors())

	freteRapidoService := freterapido.NewFreteRapidoService(viper.GetString("FRETE_RAPIDO_API_URL"))
	quoteRepository := repositories.NewQuoteRepository("quotes", m)
	quoteHandler := quotes.NewQuoteHandler(quoteRepository, freteRapidoService)

	router.Post("/quote", quoteHandler.CreateQuote)
	router.Get("/metrics", quoteHandler.GetQuoteMetrics)
	return router
}
