package routes

import (
	"time"

	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/dto"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/AbdurrahmanA/short-url/service"
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	services         *service.Services
	shortUrlDomain   string
	shortURLCacheTTL int
}

func NewShortURLRoutes(services *service.Services, shortUrlDomain string, shortURLCacheTTL int) *Routes {
	return &Routes{
		services:         services,
		shortUrlDomain:   shortUrlDomain,
		shortURLCacheTTL: shortURLCacheTTL,
	}
}

func (r *Routes) CreateNewShortURL(c *fiber.Ctx) error {
	userIP := helper.GetHerokuClintIP(c)
	body := c.Locals("body").(*request.NewURL)

	row, err := r.services.ShortURLService.Insert(body.URL, userIP, helper.CreateShortUrl)
	if err != nil {
		statusCode := fiber.StatusUnprocessableEntity
		return c.Status(statusCode).JSON(response.ErrorResponse{
			Message: err.Error(),
			Status:  statusCode,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Message: "New Short URL created successfully",
		Status:  fiber.StatusCreated,
		Data:    dto.ToUrlDTO(row, r.shortUrlDomain),
	})
}

func (r *Routes) RedirectShortURL(c *fiber.Ctx) error {
	shortURL := c.Locals("shortURL").(string)
	originalURLCache, err := r.services.RedisService.Get(shortURL)
	if err == nil {
		return c.Redirect(originalURLCache, fiber.StatusFound)
	}

	originalURL, err := r.services.ShortURLService.FindOneWithShortURL(shortURL)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
			Message: "Short URL Not Found",
			Status:  fiber.StatusNotFound,
		})
	}

	go r.services.RedisService.Set(shortURL, originalURL, time.Duration(r.shortURLCacheTTL*24*int(time.Hour)))
	return c.Redirect(originalURL, fiber.StatusFound)
}
