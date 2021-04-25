package controllers

import (
	"net/http"
	"strconv"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/errors"
	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

const WorksIDKey = "id"

type WorksController struct {
	service services.WorksService
}

//NewWorksController add /works
func NewWorksController(service services.WorksService) *WorksController {
	if service == nil {
		panic("service can't be nil")
	}

	return &WorksController{
		service: service,
	}
}

func (ctrl *WorksController) Get(c *gin.Context) {
	res, err := ctrl.service.GetAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *WorksController) FindByID(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.DSWE01), errors.Cause(err)))
		return
	}

	res, err := ctrl.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *WorksController) Post(c *gin.Context) {
	form := &beans.WorksFormBean{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(err)
		return
	}

	if err := ctrl.service.Create(c.Request.Context(), form); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *WorksController) Put(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.DSWE01), errors.Cause(err)))
		return
	}

	form := &beans.WorksFormBean{}
	if err := c.ShouldBind(form); err != nil {
		c.Error(err)
		return
	}

	if err := ctrl.service.Update(c.Request.Context(), id, form); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *WorksController) Delete(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.DSWE01), errors.Cause(err)))
		return
	}

	if err := ctrl.service.DeleteByID(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
}

func extractWorksID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param(WorksIDKey), 10, 64)
}
