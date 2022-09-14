package service

import (
	model "github.com/AbdurrahmanA/short-url/model"
	service "github.com/AbdurrahmanA/short-url/service"
)

type MockURLService struct {
	insertMethod              func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error)
	findOneWithShortURLMethod func(shortURL string) (string, error)
	findMethod                func(filter interface{}) ([]model.URL, error)
	deleteManyMethod          func(filter interface{}) error
}

func (m MockURLService) Insert(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error) {
	return m.insertMethod(url, ip, createShortUrl)
}

func (m MockURLService) FindOneWithShortURL(shortURL string) (string, error) {
	return m.findOneWithShortURLMethod(shortURL)
}

func (m MockURLService) Find(filter interface{}) ([]model.URL, error) {
	return m.findMethod(filter)
}

func (m MockURLService) DeleteMany(filter interface{}) error {
	return m.deleteManyMethod(filter)
}

func NewMockURLService(
	insertMethod func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error),
	findOneWithShortURLMethod func(shortURL string) (string, error),
	findMethod func(filter interface{}) ([]model.URL, error),
	deleteManyMethod func(filter interface{}) error,
) service.IURLService {
	mockURLService := struct{ MockURLService }{}
	mockURLService.insertMethod = insertMethod
	mockURLService.findOneWithShortURLMethod = findOneWithShortURLMethod
	mockURLService.findMethod = findMethod
	mockURLService.deleteManyMethod = deleteManyMethod
	return mockURLService
}
