package main

import (
	"github.com/AbdurrahmanA/short-url/api"
	"github.com/AbdurrahmanA/short-url/api/routes"
	env "github.com/AbdurrahmanA/short-url/pkg/Env"
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
	MongoDB = mongo.NewMongoDBClient(Env.MongoURL, Env.MongoUserName, Env.MongoPass, Env.MongoCollection, Env.MongoDBName)
	Redis = redis.NewRedisClient(Env.RedisURL, Env.RedisPass)

	shortURLRepository := repository.NewURLRepository(MongoDB)
	shortURLService := service.NewURLService(shortURLRepository, Redis, Env.URLCacheTTL, Env.URLExpirationTime)

	Services = service.RegisterServices(shortURLService)
	Routes = routes.NewShortURLRoutes(Services)
	AppPort = Env.Port

}

func main() {
	defer MongoDB.Close()
	defer Logger.Sync()
	defer Redis.Close()

	scheduler.InitExpiredScheduler(Services.ShortURLService.FindExpiredURLs, Services.ShortURLService.DeleteExpiredURLs)
	api.InitAPI(AppPort, Env.UserHourlyLimit, Routes)
}
