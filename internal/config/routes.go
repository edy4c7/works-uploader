package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/edy4c7/works-uploader/internal/controllers"
	"github.com/edy4c7/works-uploader/internal/infrastructures"
	"github.com/edy4c7/works-uploader/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const apiPath = "/api"

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	tranRnr := infrastructures.NewTransactionRunnerImpl(db)
	worksRepo := infrastructures.NewWorksRepositoryImpl(db)
	actRepo := infrastructures.NewActivitiesRepositoryImpl(db)
	userRepo := infrastructures.NewUserRepositoryImpl(db)
	uuidGen := &infrastructures.UUIDGeneratorImpl{}
	fileUploader := &infrastructures.FileUploaderImpl{}

	worksService := services.NewWorksServiceImpl(tranRnr, worksRepo, actRepo, uuidGen, fileUploader)
	worksCtrl := controllers.NewWorksController(worksService)

	actsService := services.NewActivitiesServiceImpl(actRepo)
	actsCtrl := controllers.NewActivitiesController(actsService)

	usersService := services.NewUsersServiceImpl(userRepo)
	usersCtrl := controllers.NewUsersController(usersService)

	api := r.Group(apiPath)
	v1 := api.Group("/v1")

	worksRoutes := v1.Group("/works")
	worksRoutes.GET("", worksCtrl.Get)
	worksRoutes.GET("/:id", worksCtrl.FindByID)
	worksRoutes.POST("", worksCtrl.Post)
	worksRoutes.PUT("/:id", worksCtrl.Put)
	worksRoutes.DELETE("/:id", worksCtrl.Delete)

	actsRoutes := v1.Group("/activities")
	actsRoutes.GET("", actsCtrl.Get)

	userRoutes := v1.Group("/users")
	userRoutes.PUT("", usersCtrl.Save)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	publicDir := wd + "/public"
	indexCtrl := controllers.NewIndexController(http.Dir(publicDir))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, apiPath) {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}, indexCtrl.Index)
}
