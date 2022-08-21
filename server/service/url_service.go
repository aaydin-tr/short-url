package service

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
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

func (u *URLService) Insert(url, ip string) error {
	temp := model.URL{
		OriginalURL: url,
		OwnerIP:     ip,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
	return u.repository.Insert(temp)
}
