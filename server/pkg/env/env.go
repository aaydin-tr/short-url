package env

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type ENV struct {
	MongoURL        string `env:"MONGO_URL,required"`
	MongoUserName   string `env:"MONGO_USERNAME,required"`
	MongoPass       string `env:"MONGO_PASS,required"`
	MongoDB         string `env:"MONGO_DB,required"`
	MongoCollection string `env:"MONGO_COLLECTION,required"`
	RedisURL        string `env:"REDIS_URL,required"`
	RedisPass       string `env:"REDIS_PASS,required"`
	DefaultTTLDays  int    `env:"DEFAULT_TTL_DAYS,required"`
	Port            string `env:"PORT,required"`
	LogLevel        string `env:"LOG_LEVEL,required"`
}

func ParseEnv() *ENV {
	cfg := ENV{}
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error while loading .env file: %s", err)
		os.Exit(1)
	}

	if err = env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(0)
	}

	return &cfg
}