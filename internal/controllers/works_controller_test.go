package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/common/constants"
	"github.com/edy4c7/works-uploader/internal/entities"
	myErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/mocks"
	"github.com/gin-gonic/gin"
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
	t.Run("Is valid and specify both offset and limit", func(t *testing.T) {
		const endpoint = "/?offset=%d&limit=%d"

		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		offset, limit := 1, 50

		pagination := &beans.PaginationBean{
			TotalItems: 200,
			Offset:     offset,
		}
		for _, v := range worksTestData {
			pagination.Items = append(pagination.Items, v)
		}

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx, offset, limit).Return(pagination, nil)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, offset, limit), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(pagination)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Is valid and specify only offset", func(t *testing.T) {
		const endpoint = "/?offset=%d"

		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		offset, limit := 1, 100

		pagination := &beans.PaginationBean{
			TotalItems: 200,
			Offset:     offset,
		}
		for _, v := range worksTestData {
			pagination.Items = append(pagination.Items, v)
		}

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx, offset, limit).Return(pagination, nil)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, offset), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(pagination)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Is valid and specify only limit", func(t *testing.T) {
		const endpoint = "/?limit=%d"

		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		offset, limit := 0, 50

		pagination := &beans.PaginationBean{
			TotalItems: 200,
			Offset:     offset,
		}
		for _, v := range worksTestData {
			pagination.Items = append(pagination.Items, v)
		}

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx, offset, limit).Return(pagination, nil)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, limit), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(pagination)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Is valid and not specify both offset and limit", func(t *testing.T) {
		const endpoint = "/"

		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		offset, limit := 0, 100

		pagination := &beans.PaginationBean{
			TotalItems: 200,
			Offset:     offset,
		}
		for _, v := range worksTestData {
			pagination.Items = append(pagination.Items, v)
		}

		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(ctx, offset, limit).Return(pagination, nil)
		workCtrl := NewWorksController(service)
		r.GET("/", workCtrl.Get)

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(pagination)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		const endpoint = "/"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		errExpect := errors.New("ERROR")
		service := mocks.NewMockWorksService(ctrl)
		service.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errExpect)
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
	const endpoint = "/%v"

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

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, id), nil)
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

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, "abc"), nil)
		req = req.WithContext(ctx)
		ginCtx.Request = req
		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		if err != nil {
			var appErr *myErr.ApplicationError
			if !errors.As(err.Err, &appErr) {
				assert.Fail(t, err.Err.Error())
			} else {
				assert.Equal(t, myErr.WUE01, appErr.Code())
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
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, id), nil)
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
	const endpoint = "/"

	contentType := constants.WorkType(constants.ContentTypeURL)
	title := gererateRandString(40)
	description := gererateRandString(200)
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
		expect := &entities.Work{
			ID:          12345,
			Title:       form.Title,
			Description: form.Description,
			ContentURL:  form.ContentURL,
		}
		service.EXPECT().Create(ctx, &form).Return(expect, nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		var res entities.Work
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, strings.HasSuffix(w.Header().Get("Location"), fmt.Sprintf("/%d", res.ID)))
		assert.Equal(t, *expect, res)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
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
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errExpect)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
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
		expect := &entities.Work{
			Title:       title,
			Description: description,
			ContentURL:  url,
		}
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expect, nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
	})

	t.Run("description over than 200 charcters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		description := gererateRandString(201)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
		} else {
			assert.Fail(t, "%v", errActual)
		}
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("title over than 40 characters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, r := gin.CreateTestContext(w)

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)
		title := gererateRandString(41)
		createWorksFormRequestBody(mw, contentType, title, description, url, nil, nil, 0)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestPostWorksWithFile(t *testing.T) {
	const endpoint = "/"

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
		expect := &entities.Work{
			ID:           12345,
			Title:        form.Title,
			Description:  form.Description,
			ThumbnailURL: "https://example.com/thumbnail",
			ContentURL:   "https://example.com/contenturl",
		}
		service.EXPECT().Create(ctx, &form).Return(expect, nil)
		workCtrl := NewWorksController(service)
		r.POST("/", workCtrl.Post)

		r.HandleContext(ginCtx)

		err := ginCtx.Errors.Last()
		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		var res entities.Work
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, strings.HasSuffix(w.Header().Get("Location"), fmt.Sprintf("/%d", res.ID)))
		assert.Equal(t, *expect, res)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
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
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errExpect)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
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
		service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&entities.Work{}, nil)
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
			var bre *myErr.BadRequestError
			assert.True(t, errors.As(errActual.Err, &bre))
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

func TestDeleteWorks(t *testing.T) {
	const endpoint = "/%v"

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
				assert.Equal(t, myErr.WUE01, appErr.Code())
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func gererateRandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
