package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	fileServer http.Handler
}

func NewIndexController(publicDir http.Dir) *IndexController {
	return &IndexController{
		fileServer: http.FileServer(publicDir),
	}
}

func (ctrl *IndexController) Index(c *gin.Context) {
	ctrl.fileServer.ServeHTTP(c.Writer, c.Request)
}
