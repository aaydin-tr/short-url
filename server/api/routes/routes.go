package routes

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/dto"
	"github.com/AbdurrahmanA/short-url/service"
	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	services *service.Services
}

func NewShortURLRoutes(services *service.Services) *Routes {
	return &Routes{
		services: services,
	}
}

func (r *Routes) CreateNewShortURL(c *fiber.Ctx) error {
	userIP := c.IP()
	body := c.Locals("body").(*request.NewURL)

	row, err := r.services.ShortURLService.Insert(body.URL, userIP)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Bad Request",
			Status:  fiber.StatusUnprocessableEntity,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Message: "New Short URL created successfully",
		Status:  fiber.StatusCreated,
		Data:    dto.ToUrlDTO(row),
	})
}
