package helper

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const sha256Length = 64
const availableBlockCount = 8

func CreateShortUrl(originalURL, ownerIP string) string {
	now := time.Now()
	byteArr := []byte(originalURL + now.String() + ownerIP)
	startBlock := now.Second()
	endBlock := startBlock + availableBlockCount

	if endBlock > 64 {
		startBlock = randomInt(0, sha256Length-availableBlockCount)
		endBlock = startBlock + availableBlockCount
	}

	hash := sha256.Sum256(byteArr)
	return fmt.Sprintf("%x", hash)[startBlock:endBlock]
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "url":
		return "Invalid URL"
	case "alphanum":
		return "Input should be alphanumeric"
	case "len":
		return "Input should be 8 characters long"
	}
	return fe.Error()
}

func GetHerokuClintIP(c *fiber.Ctx) string {
	xForwardedFor := c.IPs()

	if len(xForwardedFor) == 0 || xForwardedFor == nil {
		return c.IP()
	}

	lastXForwardedFor := xForwardedFor[len(xForwardedFor)-1]
	isValidIP := isValidIPAddress(lastXForwardedFor)
	if isValidIP {
		return lastXForwardedFor
	}

	return c.IP()
}

func isValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusBadRequest

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(response.ErrorResponse{
		Message: err.Error(),
		Status:  code,
	})
}

func LimiterHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(response.ErrorResponse{
		Message: "Too many attempts please try again later",
		Status:  fiber.StatusTooManyRequests,
	})
}

func NotifyShutdown(shutdown chan os.Signal) {
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}
