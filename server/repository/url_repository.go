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

func (u *URLRepository) Find(filter interface{}) ([]model.URL, error) {
	var results []model.URL
	opt := options.Find().SetProjection(bson.M{"_id": 1, "short_url": 1})

	cursor, err := u.collection.Find(u.context, filter, opt)
	if err != nil {
		return nil, err
	}

	err = cursor.All(u.context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (u *URLRepository) DeleteMany(filter interface{}) error {
	_, err := u.collection.DeleteMany(u.context, filter)
	if err != nil {
		return err
	}

	return nil
}
