package service

import (
	"context"
	"testing"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/strutil"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test{{.ModelName}}(t *testing.T) {
	t.Parallel()

	opts, err := repository.BuildNewRepositoryOptionsForTest()
	assert.NoError(t, err)
	repo := repository.New{{.ModelName}}Repository(*opts)
	svc := New{{.ModelName}}Service(repo)

	// Test Create
	suffix, err := strutil.GenerateRandomString(10)
	assert.NoError(t, err)
	userID := "testUser" + suffix
	object := &model.{{.ModelName}}{}
	object.UserID = userID

	id, err := svc.Create(context.Background(), object)
	assert.NoError(t, err)
	object.ID = id

	_, err = svc.Get(context.Background(), object.ID)
	assert.NoError(t, err)

	objects, err := svc.List(context.Background(), &model.{{.ModelName}}ListOption{
		UserID: userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 1)
	assert.Equal(t, id, objects[0].ID)

	objects, err = svc.List(context.Background(), &model.{{.ModelName}}ListOption{
		UserID: "another" + userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	err = svc.Delete(context.Background(), object.ID)
	assert.NoError(t, err)

	objects, err = svc.List(context.Background(), &model.{{.ModelName}}ListOption{
		UserID: userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	_, err = svc.Get(context.Background(), object.ID)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
