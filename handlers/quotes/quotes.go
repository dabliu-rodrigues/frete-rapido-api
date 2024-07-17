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
		utils.SendGenericError(w, http.StatusBadRequest, "At least one volume is necessary to simulate quote!")
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

	var transformedOffers = models.CreateQuoteResponse{}

	for _, offer := range resp.Dispatchers[0].Offers {
		offerData := struct {
			Name     string  `json:"name"`
			Service  string  `json:"service"`
			Deadline int     `json:"deadline"`
			Price    float64 `json:"price"`
		}{
			Name:     offer.Carrier.Name,
			Service:  offer.Service,
			Price:    offer.FinalPrice,
			Deadline: offer.DeliveryTime.Days,
		}
		transformedOffers.Carrier = append(transformedOffers.Carrier, offerData)
	}

	utils.SendResponse(w, transformedOffers)
}

func GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello /metrics!"))
}
