package service

import (
	"context"
	"testing"

	"github.com/rimoapp/repository-example/strutil"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test{{.ModelName}}(t *testing.T) {
	t.Parallel()

	repo := repository.NewRepository[*model.{{.ModelName}}](repository.NewRepositoryOption{UseInMemory: true})
	svc := New{{.ModelName}}Service(repo)

	// Test Create
	suffix, err := strutil.GenerateRandomString(10)
	assert.NoError(t, err)
	userID := "testUser" + suffix
	object := &model.{{.ModelName}}{}
	object.UserID = userID

	id, err := svc.Create(context.Background(), object, nil)
	assert.NoError(t, err)
	object.ID = id

	_, err = svc.Get(context.Background(), object.ID, nil)
	assert.NoError(t, err)

	orgID := "testOrg" + suffix
	keyValues := map[string]interface{}{"organization_id": orgID}
	err = svc.Update(context.Background(), object.ID, keyValues, nil)
	assert.NoError(t, err)

	object, err = svc.Get(context.Background(), object.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, orgID, object.OrganizationID)

	objects, err := svc.List(context.Background(), &model.BaseAssociatedListOption[*model.{{.ModelName}}]{
		UserID: userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 1)
	assert.Equal(t, id, objects[0].ID)

	objects, err = svc.List(context.Background(), &model.BaseAssociatedListOption[*model.{{.ModelName}}]{
		UserID: "another" + userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	err = svc.Delete(context.Background(), object.ID, nil)
	assert.NoError(t, err)

	objects, err = svc.List(context.Background(), &model.BaseAssociatedListOption[*model.{{.ModelName}}]{
		UserID: userID,
	})
	assert.NoError(t, err)
	assert.Len(t, objects, 0)

	_, err = svc.Get(context.Background(), object.ID, nil)
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
