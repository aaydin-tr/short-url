package repository

import (
	"context"
	"testing"

	"github.com/AbdurrahmanA/short-url/model"
	pkgmongo "github.com/AbdurrahmanA/short-url/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var mt *mtest.T
var MockURLDatas = []model.URL{
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.1", OriginalURL: "https://example.com/0", ShortURL: "fe509fd0"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.2", OriginalURL: "https://example.com/1", ShortURL: "asdasdas"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.3", OriginalURL: "https://example.com/2", ShortURL: "qweqweqw"},
	{ID: primitive.NewObjectID(), OwnerIP: "127.0.0.4", OriginalURL: "https://example.com/3", ShortURL: "fghfghfg"},
}

func setupUrlRepo(t *testing.T) func() {
	mt = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	return func() {
		mt.Close()
	}
}

func newClient(mt *mtest.T) *pkgmongo.Mongo {
	return &pkgmongo.Mongo{
		Client:         mt.Client,
		Context:        context.Background(),
		URLsCollection: mt.Coll,
	}
}

func TestInsert(t *testing.T) {
	df := setupUrlRepo(t)
	defer df()

	mt.Run("success", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		testModel := model.URL{OwnerIP: "127.0.0.0", OriginalURL: "https://example.com/0", ShortURL: "fe509fd0"}

		result, err := urlRepo.Insert(testModel.OriginalURL, testModel.OwnerIP, testModel.ShortURL)
		assert.NotNil(t, result)
		assert.Equal(t, testModel.OriginalURL, result.OriginalURL)
		assert.Equal(t, testModel.OwnerIP, result.OwnerIP)
		assert.Equal(t, testModel.ShortURL, result.ShortURL)
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		testModel := model.URL{OwnerIP: "127.0.0.0", OriginalURL: "https://example.com/0", ShortURL: "fe509fd0"}

		result, err := urlRepo.Insert(testModel.OriginalURL, testModel.OwnerIP, testModel.ShortURL)
		assert.Equal(t, (*model.URL)(nil), result)
		assert.EqualError(t, err, "write exception: write errors: [duplicate key error]")
	})

}

func TestFindOne(t *testing.T) {
	df := setupUrlRepo(t)
	defer df()

	mt.Run("success", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{"original_url", "https://example.com/0"}}))
		_, err := urlRepo.FindOne("fe509fd0")
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Message: "mongo: no documents in result",
		}))

		result, err := urlRepo.FindOne("fe509fd0")
		assert.Equal(t, "", result)
		assert.EqualError(t, err, "write command error: [{write errors: [{mongo: no documents in result}]}, {<nil>}]")
	})
}

func TestFind(t *testing.T) {
	df := setupUrlRepo(t)
	defer df()

	mt.Run("success", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{"_id", MockURLDatas[0].ID}, {"short_url", MockURLDatas[0].ShortURL}})
		getMore := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{{"_id", MockURLDatas[1].ID}, {"short_url", MockURLDatas[1].ShortURL}})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, getMore, killCursors)

		result, err := urlRepo.Find(bson.M{})
		assert.NotNil(t, result)
		assert.NoError(t, err)
		assert.Equal(t, MockURLDatas[0].ShortURL, result[0].ShortURL)
		assert.Equal(t, MockURLDatas[1].ShortURL, result[1].ShortURL)

	})

	mt.Run("error", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		_, err := urlRepo.Find(nil)
		assert.EqualError(t, err, "document is nil")

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{"_id", MockURLDatas[0].ID}, {"short_url", 0}})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		result, err := urlRepo.Find(bson.M{})
		assert.Equal(t, []model.URL(nil), result)
		assert.EqualError(t, err, "error decoding key short_url: cannot decode 32-bit integer into a string type")
	})
}

func TestDeleteMany(t *testing.T) {
	df := setupUrlRepo(t)
	defer df()

	mt.Run("success", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})
		err := urlRepo.DeleteMany(bson.M{})
		assert.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		newClient := newClient(mt)
		urlRepo := NewURLRepository(newClient)

		mt.AddMockResponses(bson.D{{"ok", 0}, {"acknowledged", false}, {"n", 0}})
		err := urlRepo.DeleteMany(bson.M{})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "command failed")
	})
}
