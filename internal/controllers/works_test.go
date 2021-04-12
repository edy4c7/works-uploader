package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	"github.com/edy4c7/darkpot-school-works/internal/mocks"
	"github.com/edy4c7/darkpot-school-works/internal/test/testutil"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const contentTypeKey string = "Content-Type"

var data = []*entities.Work{
	{
		ID:           0,
		Title:        "hoge",
		Description:  "hogehoge",
		ThumbnailURL: "https://example.com",
		ContentURL:   "https://example.com",
	},
	{
		ID:           1,
		Title:        "hoge",
		Description:  "hogehoge",
		ThumbnailURL: "https://example.com",
		ContentURL:   "https://example.com",
	},
}

const path = "/"

func TestGetWorks(t *testing.T) {
	const endpoint string = "/works/"

	t.Run("Public mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx).Return(data, nil)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, true)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req = req.WithContext(ctx)
		w := testutil.ServeHTTP(r, req)

		//公開モードでは、作品情報の取得は認証無しで可能
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(data)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Private mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx).Return(data, nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req = req.WithContext(ctx)
		w := testutil.ServeHTTP(r, req)

		assert.True(t, called)
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(data)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(gomock.Any()).Return(nil, errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, true)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestGetWorkById(t *testing.T) {
	const endpoint string = "/works/%d"

	t.Run("Public mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		id := uint64(1)
		service.EXPECT().FindByID(ctx, id).Return(data[1], nil)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, true)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, id), nil)
		req = req.WithContext(ctx)
		w := testutil.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, err, "%T %v", err, err)
		res, _ := json.Marshal(data[1])
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Private mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		called := false
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		id := uint64(1)
		service.EXPECT().FindByID(ctx, id).Return(data[1], nil)
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, id), nil)
		req = req.WithContext(ctx)
		w := testutil.ServeHTTP(r, req)

		assert.True(t, called)
		assert.Nil(t, err, "%T %v", err, err)
		res, _ := json.Marshal(data[1])
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, true)

		id := uint64(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, id), nil)
		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPostWorksWithURL(t *testing.T) {
	const endpoint string = "/works/"

	contentType := constants.WorkType(constants.ContentTypeURL)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	url := "https://example.com"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req = req.WithContext(ctx)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(ctx, &form).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
		assert.True(t, called)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, 0, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, "", url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, "", description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPostWorksWithFile(t *testing.T) {
	const endpoint string = "/works/"

	contentType := constants.WorkType(constants.ContentTypeFile)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	thumbnail := []byte{0x12, 0x34}
	content := []byte{0xab, 0xcd}

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req = req.WithContext(ctx)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(ctx, &form).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, 0, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Thumbnail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", nil, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, "", "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, "", description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPutWorksWithURL(t *testing.T) {
	const endpoint string = "/works/%d"

	targetID := uint64(1234)
	contentType := constants.WorkType(constants.ContentTypeURL)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	url := "https://example.com"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, url, nil, nil, 1)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req = req.WithContext(ctx)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(ctx, targetID, &form).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, 0, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, "", url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, "", description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPutWorksWithFile(t *testing.T) {
	const endpoint string = "/works/%d"

	targetID := uint64(1234)
	contentType := constants.WorkType(constants.ContentTypeFile)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	thumbnail := []byte{0x12, 0x34}
	content := []byte{0xab, 0xcd}

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req = req.WithContext(ctx)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(req.Context(), targetID, &form).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, 0, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", thumbnail, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Thumbnail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, description, "", nil, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, title, "", "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createFormRequestBody(mw, contentType, "", description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		service := mocks.NewMockWorksService(ctrl)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		testutil.ServeHTTP(r, req)

		if errActual != nil {
			var ve validator.ValidationErrors
			if !errors.As(errActual.Err, &ve) || len(ve) <= 0 {
				assert.Fail(t, errActual.Err.Error())
			}
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestDeleteWorks(t *testing.T) {
	const endpoint string = "/works/%d"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		targetID := uint64(1234)

		var err error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			err = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().DeleteByID(ctx, targetID).Return(nil)
		called := false
		NewWorksController(r.Group(path), service, testutil.AssertCalled(&called), false)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)
		req = req.WithContext(ctx)
		testutil.ServeHTTP(r, req)

		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1234)

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := testutil.CreateRouter(testutil.CreateErrorMiddleware(func(c *gin.Context) {
			errActual = c.Errors.Last()
		}))
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(errExpect)
		NewWorksController(r.Group(path), service, testutil.NOPHandler, false)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)
		testutil.ServeHTTP(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func createFormRequestBody(w *multipart.Writer, contentType constants.WorkType, title string, description string, url string, thumbnail []byte, content []byte, version uint) error {
	if contentType >= 0 {
		w.WriteField("type", fmt.Sprint(contentType))
	}

	if title != "" {
		w.WriteField("title", title)
	}

	if url != "" {
		w.WriteField("url", url)
	}

	if description != "" {
		w.WriteField("description", description)
	}

	if thumbnail != nil {
		pw, err := w.CreateFormFile("thumbnail", "thumbnail.png")
		if err != nil {
			return err
		}
		io.Copy(pw, bytes.NewBuffer(thumbnail))
	}

	if content != nil {
		pw, err := w.CreateFormFile("content", "content.png")
		if err != nil {
			return err
		}
		io.Copy(pw, bytes.NewBuffer(content))
	}

	if version > 0 {
		w.WriteField("version", fmt.Sprint(version))
	}

	return nil
}
