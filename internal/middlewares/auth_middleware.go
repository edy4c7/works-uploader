package middlewares

import (
	"net/http"

	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware interface {
	CheckJWT(w http.ResponseWriter, r *http.Request) error
}

const UserProperty = "user"

func NewAuthenticationMiddleware(service services.JWTAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.Authenticate(c.Writer, c.Request); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}
