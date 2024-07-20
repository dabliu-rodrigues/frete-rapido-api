package quotes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/services"
	freterapido "github.com/jsGolden/frete-rapido-api/services/frete-rapido"
	"github.com/jsGolden/frete-rapido-api/transformers"
	"github.com/jsGolden/frete-rapido-api/utils"
)

type QuoteHandler struct {
	Mongo       *services.MongoService
	FreteRapido *freterapido.Service
}

func NewQuoteHandler(Mongo *services.MongoService, FreteRapido *freterapido.Service) *QuoteHandler {
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
	documentsToInsert := make([]interface{}, len(resp.Dispatchers[0].Offers))

	for i, offer := range resp.Dispatchers[0].Offers {
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
		documentsToInsert[i] = offerData
	}

	_, err = collection.InsertMany(context.TODO(), documentsToInsert)

	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert offer: %s", err))
		return
	}

	utils.SendOKResponse(w, transformedOffers)
}

func (q *QuoteHandler) GetQuoteMetrics(w http.ResponseWriter, r *http.Request) {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$sort", bson.D{{"price", -1}}},
		},
	}

	if r.URL.Query().Has("last_quotes") {
		limit, err := strconv.ParseUint(r.URL.Query().Get("last_quotes"), 10, 64)
		if err != nil {
			utils.SendBadParamError(w, []utils.ParamError{
				{Param: "last_quotes", Message: "last_quotes param need to be a positive integer", Type: "Query param"},
			})
			return
		}
		pipeline = append(pipeline, bson.D{{"$limit", limit}})
	}

	db, err := q.Mongo.GetConnection()
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}

	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", "$name"},
		{"count", bson.D{{"$sum", 1}}},
		{"total_price", bson.D{{"$sum", "$price"}}},
		{"average_price", bson.D{{"$avg", "$price"}}},
	}}})

	pipeline = append(pipeline, bson.D{{"$sort", bson.D{{"total_price", 1}}}})
	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", nil},
		{"cheapest_quote", bson.D{{"$first", "$total_price"}}},
		{"most_expensive_quote", bson.D{{"$last", "$total_price"}}},
		{"services", bson.D{{"$push", "$$ROOT"}}},
	}}})

	collection := db.Collection("quotes")
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		utils.SendGenericError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	defer cursor.Close(context.TODO())

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
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
