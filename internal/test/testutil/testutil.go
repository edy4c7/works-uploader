package testutil

import (
	"github.com/gin-gonic/gin"
)

func NOPHandler(c *gin.Context) {}

func AssertCalled(called *bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		*called = true
	}
}
