package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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

func TestGetWorks(t *testing.T) {
	const endpoint string = "/works/"

	t.Run("Public mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		var err error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.public = true
				wco.service.EXPECT().GetAll(req.Context()).Return(data, nil)
				wco.errorMiddleware = func(c *gin.Context) {
					err = c.Errors.Last()
				}
			},
		)

		//公開モードでは、作品情報の取得は認証無しで可能
		w := executeHandler(r, req)

		assert.Nil(t, err, "%T %v", err, err)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(data)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Private mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		var err error
		called := false
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().GetAll(req.Context()).Return(data, nil)
				wco.authMiddleware = func(c *gin.Context) {
					called = true
				}
				wco.errorMiddleware = func(c *gin.Context) {
					err = c.Errors.Last()
				}
			},
		)

		w := executeHandler(r, req)

		assert.Nil(t, err, "%T %v", err, err)
		assert.True(t, called)
		assert.Equal(t, http.StatusOK, w.Code)
		res, _ := json.Marshal(data)
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().GetAll(gomock.Any()).Return(nil, errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, targetID), nil)

		var err error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.public = true
				wco.service.EXPECT().FindByID(req.Context(), targetID).Return(data[1], nil)
				wco.errorMiddleware = func(c *gin.Context) {
					err = c.Errors.Last()
				}
			},
		)

		w := executeHandler(r, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, err, "%T %v", err, err)
		res, _ := json.Marshal(data[1])
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("Private mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, targetID), nil)
		var err error
		called := false
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().FindByID(req.Context(), targetID).Return(data[1], nil)
				wco.authMiddleware = func(c *gin.Context) {
					called = true
				}
				wco.errorMiddleware = func(c *gin.Context) {
					err = c.Errors.Last()
				}
			},
		)

		w := executeHandler(r, req)
		assert.True(t, called)
		assert.Nil(t, err, "%T %v", err, err)
		res, _ := json.Marshal(data[1])
		assert.Equal(t, res, w.Body.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(endpoint, targetID), nil)
		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(req.Context(), &form).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, 0, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, "", url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, "", description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(req.Context(), &form).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, 0, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", nil, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, "", "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, "", description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(req.Context(), targetID, &form).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, 0, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to ContentURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, "", url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, "", description, url, nil, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var form beans.WorksFormBean
		if err := testutil.BindFormToObject(req, &form); err != nil {
			assert.FailNow(t, err.Error())
		}

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(req.Context(), targetID, &form).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing content type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, 0, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		executeHandler(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})

	t.Run("Missing to Content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", thumbnail, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, description, "", nil, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, title, "", "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Missing to Title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		createFormRequestBody(mw, contentType, "", description, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, targetID), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1234)
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().DeleteByID(req.Context(), targetID).Return(nil)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		assert.Nil(t, errActual, "%T %v", errActual, errActual)
	})

	t.Run("Is fail(500)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetID := uint64(1234)
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, targetID), nil)

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().DeleteByID(req.Context(), targetID).Return(errExpect)
				wco.errorMiddleware = func(c *gin.Context) {
					errActual = c.Errors.Last()
				}
			},
		)

		testutil.ExecuteHandler(r, req)

		if errActual != nil {
			assert.True(t, errors.Is(errActual.Err, errExpect), "%w", errActual.Err)
		} else {
			assert.Fail(t, "%v", errActual)
		}
	})
}

type worksControllerOptions struct {
	authMiddleware  gin.HandlerFunc
	errorMiddleware gin.HandlerFunc
	service         *mocks.MockWorksService
	public          bool
}

func expectWorksController(ctrl *gomock.Controller, options ...func(*worksControllerOptions)) *gin.Engine {
	r := gin.New()

	wco := &worksControllerOptions{
		authMiddleware: func(c *gin.Context) {},
		service:        mocks.NewMockWorksService(ctrl),
	}

	for _, opt := range options {
		opt(wco)
	}

	if wco.errorMiddleware != nil {
		r.Use(func(c *gin.Context) {
			c.Next()
			wco.errorMiddleware(c)
		})
	}

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			dr, _ := httputil.DumpRequest(c.Request, true)
			log.Printf("%q", dr)
		}
	})

	NewWorksController(r.Group("/"),
		wco.service,
		wco.authMiddleware,
		wco.public,
	)

	return r
}

func createFormRequestBody(w *multipart.Writer, contentType constants.WorkType, title string, description string, url string, thumbnail []byte, content []byte) error {
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

	return nil
}

func executeHandler(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}
