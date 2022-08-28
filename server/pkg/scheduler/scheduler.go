package scheduler

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

func InitExpiredScheduler(findExpiredURLs func() ([]model.URL, error), deleteExpiredURLs func(rows []model.URL) error, redisDelete func(key string) error) {
	expiredScheduler := gocron.NewScheduler(time.UTC)
	expiredScheduler.Every(1).Days().Do(func(
		findExpiredURLs func() ([]model.URL, error),
		deleteExpiredURLs func(rows []model.URL) error,
		redisDelete func(key string) error,
	) {
		rows, err := findExpiredURLs()
		if err != nil {
			zap.S().Error("Error while Expired URL Scheduler", err)
			return
		}

		if len(rows) == 0 {
			zap.S().Info("No expired URLs")
			return
		}

		for _, row := range rows {
			redisDelete(row.ShortURL)
		}

		err = deleteExpiredURLs(rows)
		if err != nil {
			zap.S().Error("Error while deleting expired URLs", err)
		}
	}, findExpiredURLs, deleteExpiredURLs, redisDelete)

	expiredScheduler.SingletonMode()
	expiredScheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)
	expiredScheduler.StartAsync()
}
