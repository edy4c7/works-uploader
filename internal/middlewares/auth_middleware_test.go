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
	t.Run("Is ignored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(
			IgnoreAuth(func(r *http.Request) bool {
				return r.URL.Path == "/"
			}),
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
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

	t.Run("Authorization and Authentication succeeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
			Authenticate(Path(http.MethodGet, "/"), PermitAll()),
		)

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
		middleware := NewAuthenticationMiddleware(
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
		)

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

	t.Run("Authenticator not specified", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
		)

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
	})

	t.Run("Authentication failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
			Authenticate(
				Path(http.MethodGet, "/"),
				func(r *http.Request) bool { return false },
			),
		)

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
	})

	t.Run("Specified multiple authenticator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(
			SetAuthorizer(mockJWTMiddleware.CheckJWT),
			Authenticate(Path(http.MethodGet, "/"), PermitAll()),
			Authenticate(Path(http.MethodPost, "/"), func(r *http.Request) bool { return false }),
		)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})
		r.POST("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), gomock.Any()).AnyTimes()
		c.Request = req
		r.HandleContext(c)
		assert.True(t, called)

		called = false
		req, _ = http.NewRequest(http.MethodPost, "/", nil)
		c.Request = req
		r.HandleContext(c)

		assert.False(t, called)
	})
}
