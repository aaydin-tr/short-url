package helper

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestIsValidIPAddress(t *testing.T) {
	validIPV4 := "10.40.210.253"

	isValid := isValidIPAddress(validIPV4)
	assert.True(t, isValid)

	invalidIPV4 := "1000.40.210.253"

	isValid = isValidIPAddress(invalidIPV4)
	assert.False(t, isValid)
}

func TestGetHerokuClintIP(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	t.Run("only-remote-addr", func(t *testing.T) {
		result := GetHerokuClintIP(c)

		assert.Equal(t, "0.0.0.0", result)
	})

	t.Run("not-valid-ip", func(t *testing.T) {
		c.Request().Header.Set(fiber.HeaderXForwardedFor, "asd")
		result := GetHerokuClintIP(c)

		assert.Equal(t, "0.0.0.0", result)
	})

	t.Run("valid-ip", func(t *testing.T) {
		c.Request().Header.Set(fiber.HeaderXForwardedFor, "127.0.0.1")
		result := GetHerokuClintIP(c)

		assert.Equal(t, "127.0.0.1", result)
	})

	t.Run("valid-and-multi", func(t *testing.T) {
		c.Request().Header.Set(fiber.HeaderXForwardedFor, "127.0.0.1,1.1.1.1")
		result := GetHerokuClintIP(c)

		assert.Equal(t, "1.1.1.1", result)
	})
}
