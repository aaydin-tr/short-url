package main

import (
	"github.com/AbdurrahmanA/short-url/api"
	"github.com/AbdurrahmanA/short-url/api/routes"
	"github.com/AbdurrahmanA/short-url/pkg/env"
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
	MongoDB = mongo.NewMongoDBClient(Env.MongoURL, Env.MongoCollection, Env.MongoDBName)
	Redis = redis.NewRedisClient(Env.RedisURL, Env.RedisPass)

	shortURLRepository := repository.NewURLRepository(MongoDB)
	shortURLService := service.NewURLService(shortURLRepository)

	shortURLRedisRepository := repository.NewRedisRepository(Redis)
	shortURLRedisService := service.NewRedisService(shortURLRedisRepository)

	Services = service.RegisterServices(shortURLService, shortURLRedisService)
	Routes = routes.NewShortURLRoutes(Services, Env.ShortURLDomain, Env.URLCacheTTL)
	AppPort = Env.Port

}

func main() {
	defer MongoDB.Close()
	defer Logger.Sync()
	defer Redis.Close()

	scheduler.InitExpiredScheduler(Services.ShortURLService.Find, Services.ShortURLService.DeleteMany, Services.RedisService.Delete)
	api.InitAPI(AppPort, Env.UserHourlyLimit, Routes)
}
