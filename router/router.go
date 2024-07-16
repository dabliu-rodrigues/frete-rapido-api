package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jsGolden/frete-rapido-api/handlers/quotes"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/quote", quotes.CreateQuote)
	router.Get("/metrics", quotes.GetQuoteMetrics)
	return router
}
