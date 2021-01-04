package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

//NewAuthenticationMiddleware 認証ミドルウェアを初期化する
func NewAuthenticationMiddleware(service services.AuthService) gin.HandlerFunc {
	return func(gc *gin.Context) {
		token := strings.Replace(gc.GetHeader("Authorization"), "Bearer ", "", 1)
		if err := service.VerifyToken(context.Background(), token); err != nil {
			gc.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		gc.Next()
	}
}
