package service

import (
	"github.com/AbdurrahmanA/short-url/model"
)

type URLRepo interface {
	Insert(original_url, owner_ip, short_url string) (*model.URL, error)
	FindOne(url string) (string, error)
	Find(filter interface{}) ([]model.URL, error)
	DeleteMany(filter interface{}) error
}

type URLService struct {
	repository URLRepo
}

type IURLService interface {
	Insert(url, ip string, createShortUrl CreateShortUrlFunc) (*model.URL, error)
	FindOneWithShortURL(shortURL string) (string, error)
	Find(filter interface{}) ([]model.URL, error)
	DeleteMany(filter interface{}) error
}

func NewURLService(repo URLRepo) IURLService {
	return &URLService{repository: repo}
}

type CreateShortUrlFunc func(string, string) string

func (u *URLService) Insert(url, ip string, createShortUrl CreateShortUrlFunc) (*model.URL, error) {
	shortURL := createShortUrl(url, ip)
	row, err := u.repository.Insert(url, ip, shortURL)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (u *URLService) FindOneWithShortURL(shortURL string) (string, error) {
	originalURL, err := u.repository.FindOne(shortURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (u *URLService) Find(filter interface{}) ([]model.URL, error) {
	rows, err := u.repository.Find(filter)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (u *URLService) DeleteMany(filter interface{}) error {
	err := u.repository.DeleteMany(filter)
	if err != nil {
		return err
	}
	return nil
}
