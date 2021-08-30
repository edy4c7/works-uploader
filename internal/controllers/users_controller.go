package controllers

import (
	"net/http"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/services"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service services.UsersService
}

func NewUsersController(service services.UsersService) *UsersController {
	return &UsersController{service: service}
}

func (r *UsersController) Save(c *gin.Context) {
	form := &beans.UserFormBean{}
	if err := c.ShouldBindJSON(form); err != nil {
		c.Error(err)
		return
	}

	if err := r.service.Save(c.Request.Context(), form); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
