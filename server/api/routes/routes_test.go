package routes

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	mockservice "github.com/AbdurrahmanA/short-url/mocks/service"
	"github.com/AbdurrahmanA/short-url/model"
	"github.com/AbdurrahmanA/short-url/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mockServices service.Services
var mockURLService service.IURLService
var mockRedisService service.IRedisService

func setupRedisService(
	setMethod func(key string, value interface{}, ttl time.Duration) error,
	getMethod func(key string) (string, error),
	deleteMethod func(key string) error,
) func() {

	mockRedisService = mockservice.NewMockRedisService(setMethod, getMethod, deleteMethod)
	mockServices.RedisService = mockRedisService
	return func() {
		mockRedisService = nil
	}
}

func setupURLService(
	insertMethod func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error),
	findOneWithShortURLMethod func(shortURL string) (string, error),
	findMethod func(filter interface{}) ([]model.URL, error),
	deleteManyMethod func(filter interface{}) error,
) func() {

	mockURLService = mockservice.NewMockURLService(insertMethod, findOneWithShortURLMethod, findMethod, deleteManyMethod)
	mockServices.ShortURLService = mockURLService
	return func() {
		mockURLService = nil
	}
}

func TestCreateNewShortURLWithoutError(t *testing.T) {
	dt := setupURLService(
		func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error) {
			return &model.URL{OriginalURL: url, OwnerIP: ip, ShortURL: "12345678", CreatedAt: primitive.NewDateTimeFromTime(time.Now())}, nil
		},
		nil,
		nil,
		nil,
	)
	defer dt()

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Locals("body", &request.NewURL{URL: "https://github.com/AbdurrahmanA/short-url"})
	routes := NewShortURLRoutes(&mockServices, "https://github.com/", 90)
	routes.CreateNewShortURL(c)

	res := c.Response()

	assert.Equal(t, fiber.StatusCreated, res.StatusCode())
	var response = &response.SuccessResponse{}
	err := json.Unmarshal(res.Body(), response)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/12345678", response.Data.ShortURL)
	assert.Equal(t, "https://github.com/AbdurrahmanA/short-url", response.Data.OriginalURL)

}

func TestCreateNewShortURLWithError(t *testing.T) {
	dt := setupURLService(
		func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error) {
			return nil, errors.New("something went wrong")
		},
		nil,
		nil,
		nil,
	)
	defer dt()

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Locals("body", &request.NewURL{URL: "https://github.com/AbdurrahmanA/short-url"})
	routes := NewShortURLRoutes(&mockServices, "https://github.com/", 90)
	routes.CreateNewShortURL(c)
	res := c.Response()

	assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode())
	var response = &response.ErrorResponse{}
	err := json.Unmarshal(res.Body(), response)
	assert.NoError(t, err)
	assert.Equal(t, "something went wrong", response.Message)
}

func TestRedirectShortURLWithRedisCache(t *testing.T) {
	dtUrl := setupURLService(
		func(url, ip string, createShortUrl service.CreateShortUrlFunc) (*model.URL, error) {
			return nil, nil
		},
		nil,
		nil,
		nil,
	)
	dtRedis := setupRedisService(
		func(key string, value interface{}, ttl time.Duration) error {
			return nil
		},
		func(key string) (string, error) {
			return "https://github.com/AbdurrahmanA/short-url", nil
		},
		func(key string) error {
			return nil
		},
	)
	defer dtUrl()
	defer dtRedis()

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Locals("shortURL", "12345678")
	routes := NewShortURLRoutes(&mockServices, "https://github.com/", 90)
	routes.RedirectShortURL(c)

	assert.Equal(t, fiber.StatusFound, c.Response().StatusCode())
	assert.Equal(t, "https://github.com/AbdurrahmanA/short-url", c.GetRespHeader("Location"))
}

func TestRedirectShortURLWithoutError(t *testing.T) {
	dtUrl := setupURLService(
		nil,
		func(shortURL string) (string, error) { return "https://github.com/AbdurrahmanA/short-url", nil },
		nil,
		nil,
	)
	dtRedis := setupRedisService(
		func(key string, value interface{}, ttl time.Duration) error {
			return nil
		},
		func(key string) (string, error) {
			return "", errors.New("not found")
		},
		func(key string) error {
			return nil
		},
	)
	defer dtUrl()
	defer dtRedis()

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Locals("shortURL", "12345678")
	routes := NewShortURLRoutes(&mockServices, "https://github.com/", 90)
	routes.RedirectShortURL(c)

	assert.Equal(t, fiber.StatusFound, c.Response().StatusCode())
	assert.Equal(t, "https://github.com/AbdurrahmanA/short-url", c.GetRespHeader("Location"))
}

func TestRedirectShortURLWithError(t *testing.T) {
	dtUrl := setupURLService(
		nil,
		func(shortURL string) (string, error) { return "", errors.New("not found") },
		nil,
		nil,
	)
	dtRedis := setupRedisService(
		func(key string, value interface{}, ttl time.Duration) error {
			return nil
		},
		func(key string) (string, error) {
			return "", errors.New("not found")
		},
		func(key string) error {
			return nil
		},
	)
	defer dtUrl()
	defer dtRedis()

	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Locals("shortURL", "12345678")
	routes := NewShortURLRoutes(&mockServices, "https://github.com/", 90)
	routes.RedirectShortURL(c)

	res := c.Response()

	assert.Equal(t, fiber.StatusNotFound, res.StatusCode())
	var response = &response.ErrorResponse{}
	err := json.Unmarshal(res.Body(), response)
	assert.NoError(t, err)
	assert.Equal(t, "not found", response.Message)
	assert.Equal(t, fiber.StatusNotFound, response.Status)
}
