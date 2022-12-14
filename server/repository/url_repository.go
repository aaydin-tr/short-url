package repository

import (
	"context"
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	mongodb "github.com/AbdurrahmanA/short-url/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// go:generate mockgen -source=./repository/url_repository.go -destination=./mocks/repository/mockURLRepository.go -package=repository  URLRepo
type URLRepo interface {
	Insert(original_url, owner_ip, short_url string) (*model.URL, error)
	FindOne(url string) (string, error)
	Find(filter interface{}) ([]model.URL, error)
	DeleteMany(filter interface{}) error
}

type URLRepository struct {
	context    context.Context
	collection *mongo.Collection
}

func NewURLRepository(mongo *mongodb.Mongo) *URLRepository {
	return &URLRepository{
		context:    mongo.Context,
		collection: mongo.URLsCollection,
	}
}

func (u *URLRepository) Insert(original_url, owner_ip, short_url string) (*model.URL, error) {
	newShortURL := model.URL{
		ID:          primitive.NewObjectID(),
		OriginalURL: original_url,
		OwnerIP:     owner_ip,
		ShortURL:    short_url,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
	_, err := u.collection.InsertOne(u.context, newShortURL)
	if err != nil {
		return nil, err
	}
	return &newShortURL, nil
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
