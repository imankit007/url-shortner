package inits

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinit() *mongo.Client {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI("mongodb://localhost:27017/").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	return client
}

var Client *mongo.Client = DBinit()

func OpenCollection(db *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = db.Database("url-shortner-db").Collection(collectionName)

	return collection

}
