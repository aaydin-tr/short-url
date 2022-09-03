package helper

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	"time"

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
	if net.ParseIP(ip) == nil {
		return false
	}
	return true

}
