package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	Mongo_URL string
	Database  string
	instance  *MongoInstance
}

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func NewMongoService(url, database string) *MongoService {
	return &MongoService{
		Mongo_URL: url,
		Database:  database,
	}
}

func (m *MongoService) GetConnection() (*mongo.Database, error) {
	if m.instance != nil {
		return m.instance.Db, nil
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.Mongo_URL))
	if err != nil {
		return nil, err
	}

	db := client.Database(m.Database)
	m.instance = &MongoInstance{
		client,
		db,
	}

	return db, nil
}

func (m *MongoService) Disconnect() error {
	if m.instance == nil {
		return nil
	}

	if err := m.instance.Client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
