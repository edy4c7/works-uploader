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
	"net/http/httptest"
	"testing"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/common/constants"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	myErr "github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const contentTypeKey string = "Content-Type"

var worksTestData = []*entities.Work{
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

func TestNewWorksController(t *testing.T) {
	t.Run("is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)

		assert.Same(t, service, workCtrl.service)
	})

	t.Run("service is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		assert.Panics(t, func() {
			NewWorksController(nil)
		})
	})
}

func TestGetWorks(t *testing.T) {
	const endpoint string = "/"

	t.Run("Public mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx).Return(worksTestData, nil)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		//公開モードでは、作品情報の取得は認証無しで可能
		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(worksTestData)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		errExpect := errors.New("ERROR")
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(gomock.Any()).Return(nil, errExpect)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestGetWorkById(t *testing.T) {
	t.Run("Public mode", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		id := uint64(1)
		expect := worksTestData[1]
		service.EXPECT().FindByID(ctx, id).Return(expect, nil)
		workCtrl := NewWorksController(service)
		r.GET("/:id", workCtrl.FindByID)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", id), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, err, "%T %v", err, err)
		res, _ := json.Marshal(expect)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.GET("/:id", workCtrl.FindByID)

		req, _ := http.NewRequest(http.MethodGet, "/abc", nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		if err != nil {
			var appErr *myErr.ApplicationError
			if !errors.As(err.Err, &appErr) {
				assert.Fail(t, err.Err.Error())
			} else {
				assert.Equal(t, myErr.DSWE01, appErr.Code())
			}
		} else {
			assert.Fail(t, "%v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errExpect)
		workCtrl := NewWorksController(service)
		r.GET("/:id", workCtrl.FindByID)

		id := uint64(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", id), nil)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPostWorksWithURL(t *testing.T) {
	const endpoint string = "/"

	contentType := constants.WorkType(constants.ContentTypeURL)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	url := "https://example.com"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(ctx, &form).Return(nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, 0, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, "", url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, "", description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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
	const endpoint string = "/"

	contentType := constants.WorkType(constants.ContentTypeFile)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	thumbnail := []byte{0x12, 0x34}
	content := []byte{0xab, 0xcd}

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(ctx, &form).Return(nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, 0, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", nil, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, "", "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, "", description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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
	const endpoint string = "/%d"

	targetID := uint64(1234)
	contentType := constants.WorkType(constants.ContentTypeURL)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	url := "https://example.com"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 1)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/%d", targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(ctx, targetID, &form).Return(nil)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 1)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, "/abc", buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		if err != nil {
			var appErr *myErr.ApplicationError
			if !errors.As(err.Err, &appErr) {
				assert.Fail(t, err.Err.Error())
			} else {
				assert.Equal(t, myErr.DSWE01, appErr.Code())
			}
		} else {
			assert.Fail(t, "%v", err)
		}
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, 0, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, "", url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, "", description, url, nil, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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
	const endpoint string = "/%d"

	targetID := uint64(1234)
	contentType := constants.WorkType(constants.ContentTypeFile)
	title := "foo"
	description := "aaaaaaaaaaaaaa"
	thumbnail := []byte{0x12, 0x34}
	content := []byte{0xab, 0xcd}

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(req.Context(), targetID, &form).Return(nil)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, "/abc", buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		req = req.WithContext(ctx)
		ginCtx.Request = req
		var form beans.WorksFormBean
		if err := ginCtx.ShouldBind(&form); err != nil {
			assert.FailNow(t, err.Error())
		}
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		if err != nil {
			var appErr *myErr.ApplicationError
			if !errors.As(err.Err, &appErr) {
				assert.Fail(t, err.Err.Error())
			} else {
				assert.Equal(t, myErr.DSWE01, appErr.Code())
			}
		} else {
			assert.Fail(t, "%v", err)
		}
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, 0, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", thumbnail, nil, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, description, "", nil, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, title, "", "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		createWorksFormRequestBody(mw, contentType, "", description, "", thumbnail, content, 0)
		mw.Close()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())
		ginCtx.Request = req
		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.PUT("/:id", workCtrl.Put)

		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
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
	const endpoint string = "/%d"

	t.Run("Is valid", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		targetID := uint64(1234)

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().DeleteByID(ctx, targetID).Return(nil)
		workCtrl := NewWorksController(service)
		r.DELETE("/:id", workCtrl.Delete)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		workCtrl := NewWorksController(service)
		r.DELETE("/:id", workCtrl.Delete)

		req, _ := http.NewRequest(http.MethodDelete, "/abc", nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		if err != nil {
			var appErr *myErr.ApplicationError
			if !errors.As(err.Err, &appErr) {
				assert.Fail(t, err.Err.Error())
			} else {
				assert.Equal(t, myErr.DSWE01, appErr.Code())
			}
		} else {
			assert.Fail(t, "%v", err)
		}
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1234)

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		service := mocks.NewMockWorksService(ctrl)
		errExpect := errors.New("ERROR")
		service.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(errExpect)
		workCtrl := NewWorksController(service)
		r.DELETE("/:id", workCtrl.Delete)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		errActual := ginCtx.Errors.Last()
		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func createWorksFormRequestBody(w *multipart.Writer, contentType constants.WorkType, title string, description string, url string, thumbnail []byte, content []byte, version uint) error {
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
