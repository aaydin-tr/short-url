package dto

import (
	"testing"
	"time"

	"github.com/AbdurrahmanA/short-url/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToUrlDTO(t *testing.T) {
	now := time.Now()
	id := primitive.NewObjectID()
	shortUrlDomain := "github.com/"
	model := model.URL{
		ID:          id,
		OriginalURL: "github.com/AbdurrahmanA/short-url/",
		OwnerIP:     "127.0.0.1",
		ShortURL:    "12345678",
		CreatedAt:   primitive.NewDateTimeFromTime(now),
	}

	result := ToUrlDTO(&model, shortUrlDomain)

	assert.NotNil(t, result)
	assert.Equal(t, model.OriginalURL, result.OriginalURL)
	assert.Equal(t, model.CreatedAt.Time(), result.CreatedAt)
	assert.Equal(t, shortUrlDomain+model.ShortURL, result.ShortURL)
}
