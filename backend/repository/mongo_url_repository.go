package repository

import (
	"context"

	"github.com/imankit007/url-shortner/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoURLRepository struct {
	urlCollection *mongo.Collection
}

func NewMongoURLRepository(urlCollection *mongo.Collection) URLRepository {
	return &MongoURLRepository{
		urlCollection: urlCollection,
	}
}

func (r *MongoURLRepository) Save(ctx context.Context, entry model.URLEntry) error {
	_, err := r.urlCollection.InsertOne(ctx, entry)
	return err
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

func (r *MongoURLRepository) ListAll(ctx context.Context) ([]model.URLEntry, error) {
	cursor, err := r.urlCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	urlEntries := make([]model.URLEntry, 0)
	if err := cursor.All(ctx, &urlEntries); err != nil {
		return nil, err
	}

	return urlEntries, nil
}
