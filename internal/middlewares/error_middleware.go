package middlewares

import (
	"errors"
	"net/http"

	"github.com/edy4c7/works-uploader/internal/beans"
	wuErr "github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var mapStatusCode = map[string]int{
	wuErr.WUE01: http.StatusNotFound,
	wuErr.WUE02: http.StatusForbidden,
	wuErr.WUE99: http.StatusInternalServerError,
}

func NewErrorMiddleware(messagePrinter i18n.Printer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		var ve validator.ValidationErrors
		if errors.As(err.Err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var appErr *wuErr.ApplicationError
		if errors.As(err.Err, &appErr) {
			lang := c.Request.Header.Get("Accept-Language")
			errBean := &beans.ErrorBean{
				Code:    appErr.Code(),
				Message: messagePrinter.Print(lang, appErr.Code(), appErr.MessageParams()),
			}
			if sc, ok := mapStatusCode[appErr.Code()]; ok {
				c.AbortWithStatusJSON(sc, errBean)
			}

			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
