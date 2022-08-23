package service

import (
	"errors"
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepo interface {
	Insert(data interface{}) error
	FindOne(url string) (string, error)
}

type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) *redis.StringCmd
}

type URLService struct {
	repository        URLRepo
	redisShortURLRepo RedisRepo
	redisAuthRepo     RedisRepo
	userHourlyLimit   int
	shortURLCacheTTL  int
}

func NewURLService(repo URLRepo, redisShortURLRepo, redisAuthRepo RedisRepo, userHourlyLimit int, shortURLCacheTTL int) *URLService {
	return &URLService{
		repository:        repo,
		redisShortURLRepo: redisShortURLRepo,
		redisAuthRepo:     redisAuthRepo,
		userHourlyLimit:   userHourlyLimit,
		shortURLCacheTTL:  shortURLCacheTTL,
	}
}

func (u *URLService) Insert(url, ip string) (*model.URL, error) {
	createdAt := time.Now()

	userAttemptCount, err := u.redisAuthRepo.Get(ip).Int()
	if err != nil && err != redis.Nil {
		return nil, err
	} else if err == redis.Nil {
		u.redisAuthRepo.Set(ip, "0", time.Hour)
	}

	if userAttemptCount > u.userHourlyLimit {
		return nil, errors.New("Too many attempts please try again later")
	}

	u.redisAuthRepo.Set(ip, userAttemptCount+1, redis.KeepTTL)

	newShortURL := model.URL{
		OriginalURL: url,
		OwnerIP:     ip,
		ShortURL:    helper.CreateShortUrl(url, ip, createdAt),
		CreatedAt:   primitive.NewDateTimeFromTime(createdAt),
	}

	err = u.repository.Insert(newShortURL)
	if err != nil {
		return nil, err
	}

	return &newShortURL, nil
}

func (u *URLService) Get(shortURL string) (string, error) {
	shortURLCache := u.redisShortURLRepo.Get(shortURL)
	shortURLCacheErr := shortURLCache.Err()
	if shortURLCacheErr != redis.Nil {
		return shortURLCache.Val(), nil
	}

	originalURL, err := u.repository.FindOne(shortURL)
	if err != nil {
		return "", err
	}
	u.redisShortURLRepo.Set(shortURL, originalURL, time.Duration(u.shortURLCacheTTL*24*int(time.Hour)))
	return originalURL, nil
}
