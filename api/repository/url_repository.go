package repository

import (
	"context"

	"github.com/AbdurrahmanA/short-url/model"
	mongodb "github.com/AbdurrahmanA/short-url/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepository struct {
	client     *mongo.Client
	context    context.Context
	collection *mongo.Collection
}

func NewURLRepository(mongo *mongodb.Mongo) *URLRepository {
	return &URLRepository{
		client:     mongo.Client,
		context:    mongo.Context,
		collection: mongo.URLsCollection,
	}
}

func (u *URLRepository) Insert(data interface{}) error {
	_, err := u.collection.InsertOne(u.context, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *URLRepository) FindOne(url string) (*model.URL, error) {
	var result model.URL
	filter := bson.D{{Key: "original_url", Value: url}}

	err := u.collection.FindOne(u.context, filter).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return &result, nil
}
