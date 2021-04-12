package testutil

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func CreateRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			dr, _ := httputil.DumpRequest(c.Request, true)
			log.Printf("%q", dr)
		}
	})
	r.Use(middlewares...)

	return r
}

func NOPHandler(c *gin.Context) {}

func AssertCalled(called *bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		*called = true
	}
}

func HandleError(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		handler(c)
	}
}

func BindFormToObject(req *http.Request, obj interface{}) error {
	return (&gin.Context{Request: req}).ShouldBind(obj)
}
