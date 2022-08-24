package helper

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
)

const sha256Length = 64
const availableBlockCount = 8

func CreateShortUrl(originalURL, ownerIP string, createdAt time.Time) string {
	byteArr := []byte(originalURL + createdAt.String() + ownerIP)
	startBlock := createdAt.Second()
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
