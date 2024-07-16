package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jsGolden/frete-rapido-api/handlers/quotes"
	"github.com/jsGolden/frete-rapido-api/middlewares"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.Cors())

	router.Post("/quote", quotes.CreateQuote)
	router.Get("/metrics", quotes.GetQuoteMetrics)
	return router
}
