package repository

import (
	"context"

	"github.com/AbdurrahmanA/short-url/model"
	mongodb "github.com/AbdurrahmanA/short-url/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (u *URLRepository) FindOne(url string) (string, error) {
	var result model.URL
	filter := bson.D{{Key: "short_url", Value: url}}
	opt := options.FindOne().SetProjection(bson.M{"original_url": 1})

	err := u.collection.FindOne(u.context, filter, opt).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.OriginalURL, nil
}
