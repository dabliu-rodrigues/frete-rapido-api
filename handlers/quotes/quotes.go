package quotes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/repositories"
	freterapido "github.com/jsGolden/frete-rapido-api/services/frete-rapido"
	"github.com/jsGolden/frete-rapido-api/transformers"
	"github.com/jsGolden/frete-rapido-api/utils"
)

type QuoteHandler struct {
	QuoteRepository *repositories.QuoteRepository
	FreteRapido     *freterapido.Service
}

func NewQuoteHandler(qr *repositories.QuoteRepository, fr *freterapido.Service) *QuoteHandler {
	return &QuoteHandler{
		qr,
		fr,
	}
}

func (q *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
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

	resp, err := q.FreteRapido.Quote(transformedQuote)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	var transformedOffers []*models.Quote

	for _, offer := range resp.Dispatchers[0].Offers {
		offerData := &models.Quote{
			Name:     offer.Carrier.Name,
			Service:  offer.Service,
			Price:    offer.FinalPrice,
			Deadline: offer.DeliveryTime.Days,
		}
		transformedOffers = append(transformedOffers, offerData)
	}

	_, err = q.QuoteRepository.InsertManyQuotes(transformedOffers)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	utils.SendOKResponse(w, transformedOffers)
}

func (q *QuoteHandler) GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	var limit uint64

	if r.URL.Query().Has("last_quotes") {
		lq, err := strconv.ParseUint(r.URL.Query().Get("last_quotes"), 10, 64)
		if err != nil {
			utils.SendBadParamError(w, []utils.ParamError{
				{Param: "last_quotes", Message: "last_quotes param need to be a positive integer", Type: "Query param"},
			})
			return
		}
		limit = lq
	}

	result, err := q.QuoteRepository.GetQuoteMetrics(limit)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	metrics := map[string]interface{}{
		"cheapest_quote":       result[0]["cheapest_quote"],
		"most_expensive_quote": result[0]["most_expensive_quote"],
		"services":             result[0]["services"],
	}

	utils.SendOKResponse(w, metrics)
}
