package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/edy4c7/works-uploader/internal/beans"
	"github.com/edy4c7/works-uploader/internal/common"
	"github.com/edy4c7/works-uploader/internal/errors"
	"github.com/edy4c7/works-uploader/internal/services"
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
	offset, limit, err := common.ExtractOffsetAndLimit(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if limit == -1 {
		limit = 100
	}

	res, err := ctrl.service.GetAll(c.Request.Context(), offset, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *WorksController) FindByID(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.WUE01), errors.Cause(err)))
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
	if err := c.Bind(form); err != nil {
		c.Error(errors.NewBadRequestError(err.Error(), err))
		return
	}

	res, err := ctrl.service.Create(c.Request.Context(), form)
	if err != nil {
		c.Error(err)
		return
	}

	scheme := common.GetScheme(c.Request)

	c.Header("Location", fmt.Sprintf("%s://%s%s/%d", scheme, c.Request.Host, c.FullPath(), res.ID))
	c.JSON(http.StatusCreated, res)
}

func (ctrl *WorksController) Put(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.WUE01), errors.Cause(err)))
		return
	}

	form := &beans.WorksFormBean{}
	if err := c.Bind(form); err != nil {
		c.Error(errors.NewBadRequestError(err.Error(), err))
		return
	}

	res, err := ctrl.service.Update(c.Request.Context(), id, form)
	if err != nil {
		c.Error(err)
		return
	}

	scheme := common.GetScheme(c.Request)

	locaton := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, strings.Replace(c.FullPath(), ":id", fmt.Sprint(res.ID), -1))

	c.Header("Location", locaton)
	c.JSON(http.StatusOK, res)
}

func (ctrl *WorksController) Delete(c *gin.Context) {
	id, err := extractWorksID(c)
	if err != nil {
		c.Error(errors.NewApplicationError(errors.Code(errors.WUE01), errors.Cause(err)))
		return
	}

	if err := ctrl.service.DeleteByID(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func extractWorksID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param(WorksIDKey), 10, 64)
}
