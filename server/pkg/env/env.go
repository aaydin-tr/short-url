package env

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type ENV struct {
	MongoURL          string `env:"MONGO_URL,required"`
	MongoUserName     string `env:"MONGO_USERNAME,required"`
	MongoUserPass     string `env:"MONGO_PASS,required"`
	MongoDBName       string `env:"MONGO_DB_NAME,required"`
	MongoCollection   string `env:"MONGO_COLLECTION,required"`
	RedisURL          string `env:"REDIS_URL,required"`
	RedisPass         string `env:"REDIS_PASS,required"`
	URLCacheTTL       int    `env:"URL_CACHE_TTL,required"`
	URLExpirationTime int    `env:"URL_EXPIRATION_TIME,required"`
	Port              string `env:"PORT,required"`
	LogLevel          string `env:"LOG_LEVEL,required"`
	UserHourlyLimit   int    `env:"USER_HOURLY_LIMIT"`
	ShortURLDomain    string `env:"SHORT_URL_DOMAIN"`
}

var doOnce sync.Once
var Env ENV

func ParseEnv() *ENV {
	doOnce.Do(func() {
		err := godotenv.Load()

		if err != nil && err.Error() != "open file.go: no such file or directory" {
			log.Fatalf("Error while loading .env file: %s", err)
			os.Exit(1)
		}

		if err = env.Parse(&Env); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(0)
		}

		if Env.UserHourlyLimit <= 0 {
			Env.UserHourlyLimit = 10
		}

		if Env.ShortURLDomain == "" {
			Env.ShortURLDomain = "http://localhost:" + Env.Port + "/"
		}
	})
	return &Env
}
