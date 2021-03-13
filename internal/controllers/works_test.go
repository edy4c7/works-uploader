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
		ID:                0,
		Title:             "hoge",
		Description:       "hogehoge",
		ThumbnailFileName: "https://example.com",
		ContentFileName:   "https://example.com",
	},
	{
		ID:                1,
		Title:             "hoge",
		Description:       "hogehoge",
		ThumbnailFileName: "https://example.com",
		ContentFileName:   "https://example.com",
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

func TestPostWorks(t *testing.T) {
	const endpoint string = "/works/"

	t.Run("Is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, thumbnail, content)
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
				wco.service.EXPECT().Save(req.Context(), uint64(0), &form).Return(nil)
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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
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

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		createFormRequestBody(mw, title, description, thumbnail, nil)
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

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, nil, content)
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

		title := "foo"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPost, endpoint, buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, "", description, thumbnail, content)
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

func TestPutWorks(t *testing.T) {
	const endpoint string = "/works/%d"

	t.Run("Is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		targetID := uint64(1234)
		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, thumbnail, content)
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
				wco.service.EXPECT().Save(req.Context(), targetID, &form).Return(nil)
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

		buff := new(bytes.Buffer)
		mw := multipart.NewWriter(buff)

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, uint64(1234)), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		errExpect := errors.New("ERROR")
		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpect)
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

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		createFormRequestBody(mw, title, description, thumbnail, nil)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, uint64(1234)), buff)
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

		title := "foo"
		description := "aaaaaaaaaaaaaa"
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, description, nil, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, uint64(1234)), buff)
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

		title := "foo"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, title, "", thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, uint64(1234)), buff)
		req.Header.Set(contentTypeKey, mw.FormDataContentType())

		dr, _ := httputil.DumpRequest(req, true)
		log.Printf("%q", dr)

		var errActual *gin.Error
		r := expectWorksController(ctrl,
			func(wco *worksControllerOptions) {
				wco.service.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

		description := "aaaaaaaaaaaaaa"
		thumbnail := []byte{0x12, 0x34}
		content := []byte{0xab, 0xcd}
		createFormRequestBody(mw, "", description, thumbnail, content)
		mw.Close()

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, uint64(1234)), buff)
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

func createFormRequestBody(w *multipart.Writer, title string, description string, thumbnail []byte, content []byte) error {
	if title != "" {
		w.WriteField("title", title)
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
