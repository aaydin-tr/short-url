package service

import (
	"testing"

	"github.com/AbdurrahmanA/short-url/mocks/repository"
	"github.com/AbdurrahmanA/short-url/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockURLRepo *repository.MockURLRepo
var mockService *URLService

type testCase struct {
	arg1     string
	expected interface{}
	err      interface{}
}

var MockData = []model.URL{
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.1", OriginalURL: "https://example.com/0", ShortURL: "fe509fd0"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.2", OriginalURL: "https://example.com/1", ShortURL: "asdasdas"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.3", OriginalURL: "https://example.com/2", ShortURL: "qweqweqw"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.4", OriginalURL: "https://example.com/3", ShortURL: "fghfghfg"},
}

func setup(t *testing.T) func() {
	ct := gomock.NewController(t)
	defer ct.Finish()

	mockURLRepo = repository.NewMockURLRepo(ct)
	mockService = NewURLService(mockURLRepo, 90)
	return func() {
		mockService = nil
		defer ct.Finish()
	}
}

func TestGet(t *testing.T) {
	td := setup(t)
	defer td()

	tests := []testCase{
		{arg1: MockData[0].ShortURL, expected: MockData[0].OriginalURL, err: nil},
		{arg1: "", expected: "", err: mongo.ErrNoDocuments},
	}

	for _, test := range tests {
		mockURLRepo.EXPECT().FindOne(test.arg1).Return(test.expected, test.err)
		result, err := mockService.Get(test.arg1)
		if err != nil && test.err == nil {
			t.Error(err)
		}

		if test.err != nil {
			assert.Equal(t, test.err, err)
		}

		assert.Equal(t, test.expected, result)
	}
}

func TestDeleteExpiredURLs(t *testing.T) {
	td := setup(t)
	defer td()

	var ids []primitive.ObjectID
	for _, row := range MockData {
		ids = append(ids, row.ID)
	}

	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	mockURLRepo.EXPECT().DeleteMany(deleteFilter).Return(nil)
	err := mockService.DeleteExpiredURLs(MockData)

	assert.Equal(t, nil, err)
}
