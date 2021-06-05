package middlewares

import (
	"errors"
	"net/http"
	"strings"

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

		if err := c.Errors.Last(); err != nil {
			var appErr *wuErr.ApplicationError
			if errors.As(err.Err, &appErr) {
				lang := c.Request.Header.Get("Accept-Language")
				errBean := &beans.ErrorBean{
					Code:    appErr.Code(),
					Message: messagePrinter.Print(lang, appErr.Code(), appErr.MessageParams()),
				}
				if sc, ok := mapStatusCode[appErr.Code()]; ok {
					c.AbortWithStatusJSON(sc, errBean)
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, errBean)
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
					wuErr.NewApplicationError(
						wuErr.Code(fe.Tag()),
						wuErr.Cause(fe),
						wuErr.Message("Validation Error"),
						wuErr.MessageParams(fe.Namespace(), fe.Field(), strings.Split(fe.Param(), ",")),
					),
				)
			}
		}
	}
}
