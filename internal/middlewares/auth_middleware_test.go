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
	t.Run("Is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(mockJWTMiddleware)

		r := gin.Default()
		called := false
		r.GET("/", middleware, func(c *gin.Context) {
			called = true
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().CheckJWT(gomock.Any(), req).Return(nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.True(t, called)
	})

	t.Run("Is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockJWTMiddleware := mocks.NewMockJWTMiddleware(ctrl)
		middleware := NewAuthenticationMiddleware(mockJWTMiddleware)

		r := gin.Default()
		r.GET("/", middleware, func(c *gin.Context) {
			t.FailNow()
		})

		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mockJWTMiddleware.EXPECT().
			CheckJWT(gomock.Any(), req).
			DoAndReturn(func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusUnauthorized)
				return errors.New("error")
			})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	})
}
