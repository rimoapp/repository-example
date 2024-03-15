package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/rimo-backend/internal/strutil"
	"github.com/rimoapp/rimo-backend/repository"
	"github.com/stretchr/testify/assert"

	"github.com/rimoapp/rimo-backend/model"
)

func TestHandler{{.ModelName}}(t *testing.T) {
	repo := repository.NewRepositoryForTest[*model.{{.ModelName}}](repository.NewRepositoryOptionForTest{UseInMemory: true})
	h := new{{.ModelName}}Handler(repo)

	userID := "testUser" + strutil.GenerateRandomStringForTest(10)
	objID := ""
	createRequestTest := requestTest{
		name: "update comment failed",
		prepare: func(c *gin.Context) {
			body := bytes.NewBufferString(`{}`)
			c.Request, _ = http.NewRequest("POST", "/public/{{.BasePath}}", body)
			commonPrepareContext(c, userID, false)
		},
		execute: func(c *gin.Context) {
			h.Create(c)
		},
		expectedCode: 200,
		checkResponse: func(resp []byte) {
			obj := &model.{{.ModelName}}{}
			err := json.Unmarshal(resp, obj)
			if err != nil {
				t.Fatal(errors.Wrap(err, string(resp)))
			}
			objID = obj.ID
		},
	}
	runTest(t, createRequestTest)

	assert.NotEqual(t, "", objID)

	anotherUserID := "another" + strutil.GenerateRandomStringForTest(10)
	getRequestTests := []requestTest{
		{
			name: "get_by_user",
			prepare: func(c *gin.Context) {
				c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/public/{{.BasePath}}/%s", objID), nil)
				commonPrepareContext(c, userID, false)
				c.Params = append(c.Params, gin.Param{Key: "ID", Value: objID})
			},
			execute: func(c *gin.Context) {
				h.Get(c)
			},
			expectedCode: 200,
		},
		{
			name: "get_by_another_user",
			prepare: func(c *gin.Context) {
				c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/public/{{.BasePath}}/%s", objID), nil)
				commonPrepareContext(c, anotherUserID, false)
				c.Params = append(c.Params, gin.Param{Key: "ID", Value: objID})
			},
			execute: func(c *gin.Context) {
				h.Get(c)
			},
			expectedCode: 404,
		},
		{
			name: "delete_by_another_user",
			prepare: func(c *gin.Context) {
				c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/public/{{.BasePath}}/%s", objID), nil)
				commonPrepareContext(c, anotherUserID, false)
				c.Params = append(c.Params, gin.Param{Key: "ID", Value: objID})
			},
			execute: func(c *gin.Context) {
				h.Delete(c)
			},
			expectedCode: 404,
		},
		{
			name: "list_by_user",
			prepare: func(c *gin.Context) {
				c.Request, _ = http.NewRequest("GET", "/public/{{.BasePath}}", nil)
				c.Set("uid", userID)
			},
			execute: func(c *gin.Context) {
				h.List(c)
			},
			expectedCode: 200,
			checkResponse: func(resp []byte) {
				objects := []model.{{.ModelName}}{}
				err := json.Unmarshal(resp, &objects)
				if err != nil {
					t.Fatal(errors.Wrap(err, string(resp)))
				}
				assert.Equal(t, 1, len(objects))
			},
		},
		{
			name: "list_by_another_user",
			prepare: func(c *gin.Context) {
				c.Request, _ = http.NewRequest("GET", "/public/{{.BasePath}}", nil)
				c.Set("uid", anotherUserID)
			},
			execute: func(c *gin.Context) {
				h.List(c)
			},
			expectedCode: 200,
			checkResponse: func(resp []byte) {
				objects := []model.{{.ModelName}}{}
				err := json.Unmarshal(resp, &objects)
				if err != nil {
					t.Fatal(errors.Wrap(err, string(resp)))
				}
				assert.Equal(t, 0, len(objects))
			},
		},
	}
	for _, test := range getRequestTests {
		runTest(t, test)
	}
	deleteRequestTest := requestTest{
		name: "delete_by_user",
		prepare: func(c *gin.Context) {
			c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/public/{{.BasePath}}/%s", objID), nil)
			commonPrepareContext(c, userID, false)
			c.Params = append(c.Params, gin.Param{Key: "ID", Value: objID})
		},
		execute: func(c *gin.Context) {
			h.Delete(c)
		},
		expectedCode: 200,
	}
	runTest(t, deleteRequestTest)
}
