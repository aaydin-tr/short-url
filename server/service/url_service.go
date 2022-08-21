package service

import (
	"github.com/AbdurrahmanA/short-url/model"
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

func (u *URLService) Insert(url string) error {
	temp := model.URL{OriginalURL: url}
	return u.repository.Insert(temp)
}
