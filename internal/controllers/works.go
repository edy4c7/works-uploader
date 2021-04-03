package controllers

import (
	"net/http"
	"strconv"

	"github.com/edy4c7/darkpot-school-works/internal/beans"
	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

//NewWorksController add /works
func NewWorksController(rg *gin.RouterGroup, service services.WorksService, authMiddleware gin.HandlerFunc, public bool) {
	works := rg.Group("/works")

	if authMiddleware == nil {
		panic("authMiddleware must not be null")
	}

	if public {
		works.GET("/", get(service))
		works.GET("/:id", findByID(service))
	} else {
		works.GET("/", authMiddleware, get(service))
		works.GET("/:id", authMiddleware, findByID(service))
	}

	works.POST("/", authMiddleware, post(service))
	works.PUT("/:id", authMiddleware, put(service))
	works.DELETE("/:id", authMiddleware, delete(service))
}

func get(service services.WorksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := service.GetAll(c.Request.Context())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func findByID(service services.WorksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := extractID(c)
		if err != nil {
			c.Error(err)
			return
		}

		res, err := service.FindByID(c.Request.Context(), id)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func post(service services.WorksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &beans.WorksFormBean{}
		if err := c.ShouldBind(form); err != nil {
			c.Error(err)
			return
		}

		if err := service.Create(c.Request.Context(), form); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func put(service services.WorksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &beans.WorksFormBean{}
		if err := c.ShouldBind(form); err != nil {
			c.Error(err)
			return
		}

		id, err := extractID(c)
		if err != nil {
			c.Error(err)
			return
		}

		if err := service.Update(c.Request.Context(), id, form); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func delete(service services.WorksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := extractID(c)
		if err != nil {
			c.Error(err)
			return
		}

		if err := service.DeleteByID(c.Request.Context(), id); err != nil {
			c.Error(err)
			return
		}
	}
}

func extractID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}
