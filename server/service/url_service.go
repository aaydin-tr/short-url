package service

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepo interface {
	Insert(data interface{}) error
	FindOne(url string) (string, error)
	Find(filter interface{}) ([]model.URL, error)
	DeleteMany(filter interface{}) error
}

type URLService struct {
	repository  URLRepo
	shortUrlTTL int
}

// go:generate mockgen -source=./service/url_service.go -destination=./mocks/service/mockURLService.go -package=service  IURLService
type IURLService interface {
	Insert(url, ip string) (*model.URL, error)
	Get(shortURL string) (string, error)
	FindExpiredURLs() ([]model.URL, error)
	DeleteExpiredURLs(rows []model.URL) error
}

func NewURLService(repo URLRepo, shortUrlTTL int) *URLService {
	return &URLService{repository: repo, shortUrlTTL: shortUrlTTL}
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
	originalURL, err := u.repository.FindOne(shortURL)
	if err != nil {
		return "", err
	}
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
		ids = append(ids, row.ID)
	}

	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	err := u.repository.DeleteMany(deleteFilter)
	if err != nil {
		return err
	}
	return nil
}
