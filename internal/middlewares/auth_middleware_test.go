package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edy4c7/darkpot-school-works/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	t.Run("Is skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthorizationMiddleware(
			mockJWTMiddleware,
			SkipAuthorization(func(r *http.Request) bool {
				return r.URL.Path == "/"
			}),
		)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		c.Request = req
		r.HandleContext(c)

		assert.True(t, called)
	})

	t.Run("Authorization succeeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthorizationMiddleware(mockJWTMiddleware)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), req).Return(nil)
		c.Request = req
		r.HandleContext(c)

		assert.True(t, called)
	})

	t.Run("Authentication succeeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(func(r *http.Request) bool { return true })

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), req).Return(nil)
		c.Request = req
		r.HandleContext(c)

		assert.True(t, called)
	})

	t.Run("Authorization failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthorizationMiddleware(mockJWTMiddleware)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), gomock.Any()).Return(errors.New("error"))
		c.Request = req
		r.HandleContext(c)

		assert.False(t, called)
	})

	t.Run("Authentication failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(func(r *http.Request) bool { return false })

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), gomock.Any()).Return(nil)
		c.Request = req
		r.HandleContext(c)

		assert.False(t, called)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
