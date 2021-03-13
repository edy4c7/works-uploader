package middlewares

import (
	"errors"
	"net/http"
	"strings"

	myErr "github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewErrorMiddleware(messageLoader tools.MessageLoader) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			var appErr *myErr.ApplicationError
			if errors.As(err.Err, appErr) {
				wrapedErr := myErr.NewApplicationError(
					myErr.Message(messageLoader.LoadMessage(appErr.Code(), "ja-JP", appErr.MessageParams())),
					myErr.Code(appErr.Code()),
					myErr.Cause(appErr),
				)

				switch appErr.Code() {
				case myErr.DSWE01:
					c.AbortWithError(http.StatusNotFound, wrapedErr)
				case myErr.DSWE99:
					c.AbortWithError(http.StatusInternalServerError, wrapedErr)
				}
			}
		}
	}
}

func NewValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if err := c.Errors.Last(); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err.Err, &ve) && len(ve) > 0 {
				fe := ve[0]

				c.AbortWithError(http.StatusBadRequest,
					myErr.NewApplicationError(
						myErr.Code(fe.Tag()),
						myErr.Cause(fe),
						myErr.Message("Validation Error"),
						myErr.MessageParams(fe.Namespace(), fe.Field(), strings.Split(fe.Param(), ",")),
					),
				)
			}
		}
	}
}
