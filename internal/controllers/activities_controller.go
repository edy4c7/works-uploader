package controllers

import (
	"net/http"

	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

const UserKey = "user"

type ActivitiesController struct {
	service services.ActivitiesService
}

//NewActivitiesController add /activities
func NewActivitiesController(service services.ActivitiesService) *ActivitiesController {
	if service == nil {
		panic("service can't be nil")
	}
	return &ActivitiesController{
		service: service,
	}
}

func (ctrl *ActivitiesController) Get(c *gin.Context) {
	user := c.Query(UserKey)

	if user == "" {
		res, err := ctrl.service.GetAll(c.Request.Context())
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, res)
	} else {
		res, err := ctrl.service.FindByUserID(c.Request.Context(), user)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
