package quotes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/services"
	"github.com/jsGolden/frete-rapido-api/transformers"
	"github.com/jsGolden/frete-rapido-api/utils"
)

func CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quoteRequest models.CreateQuoteRequest
	err := render.DecodeJSON(r.Body, &quoteRequest)
	if err != nil {
		utils.SendGenericError(w, http.StatusBadRequest, "Malformed JSON")
		return
	}

	badParams := utils.Validator(quoteRequest)
	if badParams != nil {
		utils.SendBadParamError(w, badParams)
		return
	}

	if len(quoteRequest.Volumes) <= 0 {
		utils.SendGenericError(w, http.StatusBadRequest, "At least one volume is necessary to quote!")
		return
	}

	transformedQuote, err := transformers.TransformQuoteToFreteRapido(&quoteRequest)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	freteRapidoService := services.NewFreteRapidoService("https://sp.freterapido.com/api/v3")

	resp, err := freteRapidoService.Quote(transformedQuote)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	utils.SendResponse(w, resp)
}

func GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello /metrics!"))
}
