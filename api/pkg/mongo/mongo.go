package mongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var doOnce sync.Once
var client *mongo.Client

type Mongo struct {
	Client  *mongo.Client
	Context context.Context
}

func NewConnection(url string) *Mongo {
	context := context.Background()

	doOnce.Do(func() {
		cli, err := mongo.Connect(context, options.Client().ApplyURI(url))
		if err != nil {
			panic(err)
		}
		err = cli.Ping(context, nil)
		if err != nil {
			panic(err)
		}
		zap.S().Info("MongoDB connected successfully")
		client = cli
	})

	return &Mongo{
		Client:  client,
		Context: context,
	}
}

func (m *Mongo) Close() {
	err := m.Client.Disconnect(m.Context)
	if err != nil {
		zap.S().Error("Error while disconnecting from MongoDB", err)
	}
	zap.S().Info("MongoDB disconnected successfully")
}
