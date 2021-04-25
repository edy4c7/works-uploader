package config

import (
	"github.com/edy4c7/darkpot-school-works/internal/controllers"
	"github.com/edy4c7/darkpot-school-works/internal/infrastructures"
	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	tranRnr := infrastructures.NewTransactionRunnerImpl(db)
	worksRepo := infrastructures.NewWorksRepositoryImpl(db)
	actRepo := infrastructures.NewActivitiesRepositoryImpl(db)
	uuidGen := &infrastructures.UUIDGeneratorImpl{}
	fileUploader := &infrastructures.FileUploaderImpl{}

	worksService := services.NewWorksServiceImpl(tranRnr, worksRepo, actRepo, uuidGen, fileUploader)
	worksCtrl := controllers.NewWorksController(worksService)

	actsService := services.NewActivitiesServiceImpl(actRepo)
	actsCtrl := controllers.NewActivitiesController(actsService)

	v1 := r.Group("/v1")

	worksRoutes := v1.Group("/works")
	worksRoutes.GET("/", worksCtrl.Get)
	worksRoutes.GET("/:id", worksCtrl.FindByID)
	worksRoutes.POST("/", worksCtrl.Post)
	worksRoutes.PUT("/:id", worksCtrl.Put)
	worksRoutes.DELETE("/:id", worksCtrl.Delete)

	actsRoutes := v1.Group("/activities")
	actsRoutes.GET("/", actsCtrl.Get)
}
