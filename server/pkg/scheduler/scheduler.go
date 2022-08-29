package scheduler

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/pkg/env"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var ShortUrlTTL = &env.Env.URLExpirationTime

func InitExpiredScheduler(find func(filter interface{}) ([]model.URL, error), deleteMany func(filter interface{}) error, redisDelete func(key string) error) {
	expiredScheduler := gocron.NewScheduler(time.UTC)
	expiredScheduler.Every(1).Days().Do(func(
		find func(filter interface{}) ([]model.URL, error),
		deleteMany func(filter interface{}) error,
		redisDelete func(key string) error,
	) {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, 0, -(*ShortUrlTTL))
		filter := bson.M{"created_at": bson.M{"$lt": primitive.NewDateTimeFromTime(oneMonthAgo)}}

		rows, err := find(filter)
		if err != nil {
			zap.S().Error("Error while Expired URL Scheduler", err)
			return
		}

		if len(rows) == 0 {
			zap.S().Info("No expired URLs")
			return
		}

		var ids []primitive.ObjectID
		for _, row := range rows {
			ids = append(ids, row.ID)
			redisDelete(row.ShortURL)
		}

		deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
		err = deleteMany(deleteFilter)

		if err != nil {
			zap.S().Error("Error while deleting expired URLs", err)
		}
	}, find, deleteMany, redisDelete)

	expiredScheduler.SingletonMode()
	expiredScheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)
	expiredScheduler.StartAsync()
}
