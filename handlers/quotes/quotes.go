package quotes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/services"
	"github.com/jsGolden/frete-rapido-api/transformers"
	"github.com/jsGolden/frete-rapido-api/utils"
)

type QuoteHandler struct {
	Mongo       *services.MongoService
	FreteRapido *services.FreteRapidoService
}

func NewQuoteHandler(Mongo *services.MongoService, FreteRapido *services.FreteRapidoService) *QuoteHandler {
	return &QuoteHandler{
		Mongo,
		FreteRapido,
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

	db, err := q.Mongo.GetConnection()
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	collection := db.Collection("quotes")

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

		_, err := collection.InsertOne(context.TODO(), bson.M{
			"name":     offer.Carrier.Name,
			"nervice":  offer.Service,
			"price":    offer.FinalPrice,
			"deadline": offer.DeliveryTime.Days,
		})

		if err != nil {
			utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert offer: %s", err))
			return
		}
	}

	utils.SendOKResponse(w, transformedOffers)
}

func GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello /metrics!"))
}
