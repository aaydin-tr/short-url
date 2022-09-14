package main

import (
	"os"

	"github.com/AbdurrahmanA/short-url/api"
	"github.com/AbdurrahmanA/short-url/api/routes"
	"github.com/AbdurrahmanA/short-url/pkg/env"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/AbdurrahmanA/short-url/pkg/logger"
	"github.com/AbdurrahmanA/short-url/pkg/mongo"
	"github.com/AbdurrahmanA/short-url/pkg/redis"
	"github.com/AbdurrahmanA/short-url/pkg/scheduler"
	"github.com/AbdurrahmanA/short-url/repository"
	"github.com/AbdurrahmanA/short-url/service"
	"go.uber.org/zap"
)

var AppPort string
var Env *env.ENV
var Logger *zap.Logger

var Routes *routes.Routes
var MongoDB *mongo.Mongo
var Redis *redis.Redis
var Services *service.Services

func init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	MongoDB = mongo.NewMongoDBClient(Env.MongoURL, Env.MongoUserName, Env.MongoUserPass, Env.MongoCollection, Env.MongoDBName)
	Redis = redis.NewRedisClient(Env.RedisURL, Env.RedisPass)

	shortURLRepository := repository.NewURLRepository(MongoDB)
	shortURLService := service.NewURLService(shortURLRepository)

	shortURLRedisRepository := repository.NewRedisRepository(Redis)
	shortURLRedisService := service.NewRedisService(shortURLRedisRepository)

	Services = &service.Services{ShortURLService: shortURLService, RedisService: shortURLRedisService}
	Routes = routes.NewShortURLRoutes(Services, Env.ShortURLDomain, Env.URLCacheTTL)
	AppPort = Env.Port

}

func main() {
	defer MongoDB.Close()
	defer Logger.Sync()
	defer Redis.Close()

	shutdown := make(chan os.Signal, 1)
	serverShutdown := make(chan bool)
	go helper.NotifyShutdown(shutdown)

	scheduler.InitExpiredScheduler(Services.ShortURLService.Find, Services.ShortURLService.DeleteMany, Services.RedisService.Delete)
	api := api.NewApi(AppPort, Env.UserHourlyLimit, Routes, helper.ErrorHandler, helper.LimiterHandler)

	go func() {
		<-shutdown
		Logger.Info("Gracefully shutting down...")
		err := api.Shutdown()
		if err != nil {
			Logger.Sugar().Error("Error while Shutdown Server %s", err.Error())
		}
		serverShutdown <- true
	}()

	api.InitAPI()
	<-serverShutdown
}
