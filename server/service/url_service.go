package service

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepo interface {
	Insert(data interface{}) error
	FindOne(url string) (*model.URL, error)
}

type URLService struct {
	repository URLRepo
}

func NewURLService(repo URLRepo) *URLService {
	return &URLService{repository: repo}
}

func (u *URLService) Insert(url, ip string) (*model.URL, error) {
	createdAt := time.Now()
	newShortURL := model.URL{
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
