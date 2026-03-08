package infrastructure

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName      = "url-shortner-db"
	URLCollectionName = "url-collection"
)

var MongoClient *mongo.Client

func NewMongoClient() (*mongo.Client, error) {
	if MongoClient != nil {
		return MongoClient, nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI("mongodb://localhost:27017/").
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	MongoClient = client
	return client, nil
}

func NewURLCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(DatabaseName).Collection(URLCollectionName)
}
