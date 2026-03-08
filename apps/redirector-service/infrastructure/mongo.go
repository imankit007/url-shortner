package infrastructure

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName      = "url-shortner-db"
	URLCollectionName = "url-collection"
)

func NewMongoClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI("mongodb://localhost:27017/").
		SetServerAPIOptions(serverAPI)

	return mongo.Connect(opts)
}

func NewURLCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(DatabaseName).Collection(URLCollectionName)
}
