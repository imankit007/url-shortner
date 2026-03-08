package repository

import (
	"context"

	"github.com/imankit007/url-shortner-redirector-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoURLRepository struct {
	urlCollection *mongo.Collection
}

func NewMongoURLRepository(urlCollection *mongo.Collection) *MongoURLRepository {
	return &MongoURLRepository{
		urlCollection: urlCollection,
	}
}

func (r *MongoURLRepository) FindByCode(ctx context.Context, code int64) (model.URLEntry, error) {
	var entry model.URLEntry
	err := r.urlCollection.FindOne(ctx, bson.M{"code": code}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.URLEntry{}, ErrURLNotFound
		}
		return model.URLEntry{}, err
	}

	return entry, nil
}
