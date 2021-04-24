package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddActivitiesRoutes add /activities
func AddActivitiesRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"massage": "pong",
		})
	})
}
