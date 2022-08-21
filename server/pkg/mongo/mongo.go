package mongo

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var doOnce sync.Once
var client *mongo.Client
var collection *mongo.Collection

type Mongo struct {
	Client         *mongo.Client
	Context        context.Context
	URLsCollection *mongo.Collection
}

func NewConnection(url, username, pass, collectionName, database string) *Mongo {
	context := context.Background()

	doOnce.Do(func() {
		url = fmt.Sprintf("mongodb://%s:%s@%s", username, pass, url)
		cli, err := mongo.Connect(context, options.Client().ApplyURI(url))
		if err != nil {
			panic(err)
		}
		err = cli.Ping(context, nil)
		if err != nil {
			panic(err)
		}
		zap.S().Info("MongoDB connected successfully")

		database := cli.Database(database)
		isExists := IsCollectionExists(database, context, collectionName)

		if !isExists {
			err = database.CreateCollection(context, collectionName)
			if err != nil {
				zap.S().Error(err)
			}
		}

		collection = database.Collection(collectionName)
		client = cli
	})

	return &Mongo{
		Client:         client,
		Context:        context,
		URLsCollection: collection,
	}
}

func (m *Mongo) Close() {
	err := m.Client.Disconnect(m.Context)
	if err != nil {
		zap.S().Error("Error while disconnecting from MongoDB", err)
	}
	zap.S().Info("MongoDB disconnected successfully")
}

func IsCollectionExists(database *mongo.Database, context context.Context, collectionName string) bool {
	collections, err := database.ListCollectionNames(context, bson.D{})
	if err != nil && err != mongo.ErrNilDocument {
		zap.S().Error("Error while listing collections", err)
		return false
	}

	for _, collection := range collections {
		if collection == collectionName {
			return true
		}
	}
	return false
}
