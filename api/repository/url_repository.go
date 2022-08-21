package repository

import (
	"context"

	mongodb "github.com/AbdurrahmanA/short-url/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepository struct {
	client     *mongo.Client
	context    context.Context
	collection *mongo.Collection
}

func NewURLRepository(mongo *mongodb.Mongo, collectionName string) *URLRepository {
	return &URLRepository{
		client:     mongo.Client,
		context:    mongo.Context,
		collection: mongo.URLsCollection,
	}
}

func (u *URLRepository) Insert(url string) error {
	_, err := u.collection.InsertOne(u.context, url)
	if err != nil {
		return err
	}
	return nil
}

func (u *URLRepository) FindOne(url string) (string, error) {
	var result string
	filter := bson.D{{Key: "original_url", Value: url}}

	err := u.collection.FindOne(u.context, filter).Decode(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
