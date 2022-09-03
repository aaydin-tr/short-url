package middleware

import (
	"encoding/json"
	"testing"

	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestCreateNewShortURLValidation(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	t.Run("invalid-body", func(t *testing.T) {
		err := CreateNewShortURLValidation(c)
		assert.Equal(t, fiber.ErrUnprocessableEntity, err)
	})

	t.Run("invalid-url", func(t *testing.T) {
		body := request.NewURL{URL: "asdasd"}
		bodyBytes, _ := json.Marshal(body)
		c.Request().Header.Set("Content-Type", "application/json")
		c.Request().SetBody(bodyBytes)

		CreateNewShortURLValidation(c)
		assert.Equal(t, fiber.StatusUnprocessableEntity, c.Response().StatusCode())
	})

}
