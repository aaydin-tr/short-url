package service

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepo interface {
	Insert(data interface{}) error
	FindOne(url string) (string, error)
	Find(filter interface{}) ([]model.URL, error)
	DeleteMany(filter interface{}) error
}

type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) *redis.StringCmd
	Delete(key string) error
}

type URLService struct {
	repository        URLRepo
	redisShortURLRepo RedisRepo
	shortURLCacheTTL  int
	shortUrlTTL       int
}

func NewURLService(repo URLRepo, redisShortURLRepo RedisRepo, shortURLCacheTTL, shortUrlTTL int) *URLService {
	return &URLService{
		repository:        repo,
		redisShortURLRepo: redisShortURLRepo,
		shortURLCacheTTL:  shortURLCacheTTL,
		shortUrlTTL:       shortUrlTTL,
	}
}

func (u *URLService) Insert(url, ip string) (*model.URL, error) {
	createdAt := time.Now()
	newShortURL := model.URL{
		ID:          primitive.NewObjectID(),
		OriginalURL: url,
		OwnerIP:     ip,
		ShortURL:    helper.CreateShortUrl(url, ip, createdAt),
		CreatedAt:   primitive.NewDateTimeFromTime(createdAt),
	}

	err := u.repository.Insert(newShortURL)
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

func (u *URLService) FindExpiredURLs() ([]model.URL, error) {
	now := time.Now()
	oneMonthAgo := now.AddDate(0, 0, -(u.shortUrlTTL))
	filter := bson.M{"created_at": bson.M{"$lt": primitive.NewDateTimeFromTime(oneMonthAgo)}}

	rows, err := u.repository.Find(filter)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (u *URLService) DeleteExpiredURLs(rows []model.URL) error {
	var ids []primitive.ObjectID
	for _, row := range rows {
		err := u.redisShortURLRepo.Delete(row.ShortURL)
		if err != nil {
			return err
		}
		ids = append(ids, row.ID)
	}

	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	err := u.repository.DeleteMany(deleteFilter)
	if err != nil {
		return err
	}
	return nil
}
