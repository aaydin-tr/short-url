package service

import (
	"errors"
	"testing"
	"time"

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

// TODO: Add more test case
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
	mockService = NewURLService(mockURLRepo)
	return func() {
		mockService = nil
		defer ct.Finish()
	}
}

func TestFindOneWithShortURL(t *testing.T) {
	td := setup(t)
	defer td()

	tests := []testCase{
		{arg1: MockData[0].ShortURL, expected: MockData[0].OriginalURL, err: nil},
		{arg1: "", expected: "", err: mongo.ErrNoDocuments},
	}

	for _, test := range tests {
		mockURLRepo.EXPECT().FindOne(test.arg1).Return(test.expected, test.err)
		result, err := mockService.FindOneWithShortURL(test.arg1)
		if err != nil && test.err == nil {
			t.Error(err)
		}

		if test.err != nil {
			assert.Equal(t, test.err, err)
		}

		assert.Equal(t, test.expected, result)
	}
}

func TestDeleteMany(t *testing.T) {
	td := setup(t)
	defer td()

	var ids []primitive.ObjectID
	for _, row := range MockData {
		ids = append(ids, row.ID)
	}

	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	mockURLRepo.EXPECT().DeleteMany(deleteFilter).Return(nil)
	err := mockService.DeleteMany(deleteFilter)

	assert.Equal(t, nil, err)
}

func TestDeleteManyWithError(t *testing.T) {
	td := setup(t)
	defer td()

	var ids []primitive.ObjectID

	deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
	mockURLRepo.EXPECT().DeleteMany(deleteFilter).DoAndReturn(func(filter interface{}) error {
		return errors.New("something went wrong")
	})
	err := mockService.DeleteMany(deleteFilter)
	assert.Equal(t, errors.New("something went wrong"), err)
}

func TestInsert(t *testing.T) {
	td := setup(t)
	defer td()

	testCaseObjectID, _ := primitive.ObjectIDFromHex("630a226b6303e67f7a003a43")
	time := time.Now()

	tests := []struct {
		Model              model.URL
		CreateShortUrlFunc func(string, string) string
		Err                interface{}
	}{
		{
			Model: model.URL{
				ID:          testCaseObjectID,
				OwnerIP:     "https://example.com/0",
				OriginalURL: "127.0.0.1",
				ShortURL:    "12345678",
				CreatedAt:   primitive.NewDateTimeFromTime(time),
			},
			CreateShortUrlFunc: func(original_url, owner_ip string) string {
				return "12345678"
			},
			Err: nil,
		},
		{
			Model: model.URL{
				ID:          testCaseObjectID,
				OwnerIP:     "https://example.com/0",
				OriginalURL: "127.0.0.1",
				ShortURL:    "12345678",
				CreatedAt:   primitive.NewDateTimeFromTime(time),
			},
			CreateShortUrlFunc: func(original_url, owner_ip string) string {
				return "12345678"
			},
			Err: errors.New("something went wrong"),
		},
	}

	for _, test := range tests {
		mockURLRepo.EXPECT().Insert(test.Model.OriginalURL, test.Model.OwnerIP, test.Model.ShortURL).Return(&test.Model, test.Err)
		result, err := mockService.Insert(test.Model.OriginalURL, test.Model.OwnerIP, test.CreateShortUrlFunc)

		if err != nil && test.Err == nil {
			t.Error(err)
		}

		if test.Err != nil {
			assert.Equal(t, test.Err, err)
		} else {
			assert.Equal(t, &test.Model, result)
		}
	}

}

func TestFind(t *testing.T) {
	td := setup(t)
	defer td()

	var ids []primitive.ObjectID
	for _, row := range MockData {
		ids = append(ids, row.ID)
	}

	findFilter := bson.M{"_id": bson.M{"$in": ids}}
	mockURLRepo.EXPECT().Find(findFilter).Return(MockData, nil)
	result, err := mockService.Find(findFilter)

	assert.Equal(t, nil, err)
	assert.Equal(t, MockData, result)
}

func TestFindWithErr(t *testing.T) {
	td := setup(t)
	defer td()

	var ids []primitive.ObjectID

	findFilter := bson.M{"_id": bson.M{"$in": ids}}
	mockURLRepo.EXPECT().Find(findFilter).DoAndReturn(func(filter interface{}) ([]model.URL, error) {
		return nil, errors.New("something went wrong")
	})
	_, err := mockService.Find(findFilter)

	assert.Equal(t, errors.New("something went wrong"), err)
}
