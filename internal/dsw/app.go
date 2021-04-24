package dsw

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/edy4c7/darkpot-school-works/internal/controllers"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	"github.com/edy4c7/darkpot-school-works/internal/infrastructures"
	"github.com/edy4c7/darkpot-school-works/internal/middlewares"
	"github.com/edy4c7/darkpot-school-works/internal/services"
	"github.com/gin-gonic/gin"
)

//Run run app
func Run() {
	r := gin.Default()

	r.Use(middlewares.NewValidationErrorHandler())

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SCHEMA"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(entities.Work{}, entities.Activity{})

	jwtMiddleware := middlewares.NewJWTMiddleware(os.Getenv("AUTH0_AUDIENCE"), os.Getenv("AUTH0_DOMAIN"))
	authorizationMiddleware := middlewares.NewAuthorizationMiddleware(
		jwtMiddleware, middlewares.SkipAuthorization(func(r *http.Request) bool {
			match, _ := path.Match("v*/works/*", r.URL.Path)
			return r.Method == http.MethodGet && match
		}),
	)

	r.Use(authorizationMiddleware)

	tranRnr := infrastructures.NewTransactionRunnerImpl(db)

	v1 := r.Group("/v1")
	controllers.NewWorksController(v1.Group("/works"),
		services.NewWorksServiceImpl(
			tranRnr,
			infrastructures.NewWorksRepositoryImpl(db),
			infrastructures.NewActivitiesRepositoryImpl(db),
			infrastructures.DefaultUUIDGenerator,
			&infrastructures.FileUploaderImpl{},
		))

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}
	r.Run(fmt.Sprintf(":%d", port))
}
