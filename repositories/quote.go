package repositories

import (
	"context"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuoteRepository struct {
	CollectionName string
	Mongo          *services.MongoService
}

type MetricsResponse struct {
	CheapestQuote      float64 `bson:"cheapest_quote"`
	MostExpensiveQuote float64 `bson:"most_expensive_quote"`
	Services           []struct {
		AveragePrice float64 `bson:"average_price"`
		Carrier      string  `bson:"carrier"`
		Count        int     `bson:"count"`
		TotalPrice   float64 `bson:"total_price"`
	} `bson:"services"`
}

func NewQuoteRepository(collectionName string, m *services.MongoService) *QuoteRepository {
	return &QuoteRepository{
		collectionName,
		m,
	}
}

func (qr *QuoteRepository) getCollection() (*mongo.Collection, error) {
	db, err := qr.Mongo.GetConnection()
	if err != nil {
		return nil, err
	}

	return db.Collection(qr.CollectionName), nil
}

func (qr *QuoteRepository) InsertManyQuotes(q []*models.Quote) (*mongo.InsertManyResult, error) {
	documents := make([]interface{}, len(q))
	for i, carrier := range q {
		documents[i] = carrier
	}

	c, err := qr.getCollection()
	if err != nil {
		return nil, err
	}

	result, err := c.InsertMany(context.TODO(), documents)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (qr *QuoteRepository) GetQuoteMetrics(limit uint64) (*MetricsResponse, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$sort", bson.D{{"price", -1}}},
		},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.D{{"$limit", limit}})
	}

	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", "$name"},
		{"count", bson.D{{"$sum", 1}}},
		{"total_price", bson.D{{"$sum", "$price"}}},
		{"average_price", bson.D{{"$avg", "$price"}}},
	}}})

	pipeline = append(pipeline, bson.D{{"$project", bson.D{
		{"carrier", "$_id"},
		{"count", 1},
		{"total_price", 1},
		{"average_price", 1},
		{"_id", 0},
	}}})

	pipeline = append(pipeline, bson.D{{"$sort", bson.D{{"total_price", 1}}}})

	pipeline = append(pipeline, bson.D{{"$group", bson.D{
		{"_id", nil},
		{"cheapest_quote", bson.D{{"$first", "$total_price"}}},
		{"most_expensive_quote", bson.D{{"$last", "$total_price"}}},
		{"services", bson.D{{"$push", "$$ROOT"}}},
	}}})

	c, err := qr.getCollection()
	if err != nil {
		return nil, err
	}

	cursor, err := c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var result []*MetricsResponse
	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result[0], nil
}
