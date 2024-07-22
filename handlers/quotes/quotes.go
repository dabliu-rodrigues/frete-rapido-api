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

// @Summary Create quote
// @Description Uses frete-rapido API to simulate quote and save response at database
// @Tags Quotes
// @Accept json
// @Produce json
// @Param request body models.CreateQuoteRequest true "Request body"
// @Success 200 {object} models.CreateQuoteResponse
// @Failure 400 {object} utils.BadParamError
// @Failure 500 {object} utils.ErrorResponse
// @Router /quote [post]
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

// @Summary Get quote metrics
// @Description Use stored quote to generate general metrics
// @Tags Quotes
// @Accept json
// @Produce json
// @Param        last_quotes    query     int  false  "limit carriers"
// @Success 200 {object} models.MetricsResponse
// @Failure 400 {object} utils.BadParamError
// @Failure 500 {object} utils.ErrorResponse
// @Router /metrics [get]
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

	metrics := models.MetricsResponse{
		CheapestQuote:      result.CheapestQuote,
		MostExpensiveQuote: result.CheapestQuote,
		Services: []struct {
			AveragePrice float64 "json:\"average_price\""
			Carrier      string  "json:\"carrier\""
			Count        int     "json:\"count\""
			TotalPrice   float64 "json:\"total_price\""
		}(result.Services),
	}

	utils.SendOKResponse(w, metrics)
}
