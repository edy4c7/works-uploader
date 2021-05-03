package dsw

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/edy4c7/darkpot-school-works/internal/config"
	"github.com/edy4c7/darkpot-school-works/internal/entities"
	"github.com/edy4c7/darkpot-school-works/internal/middlewares"
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

	jwtMiddleware := middlewares.NewJWTMiddleware(os.Getenv("AUTH0_AUDIENCE"), os.Getenv("AUTH0_ISSUER"))
	authorizationMiddleware := middlewares.NewAuthorizationMiddleware(
		jwtMiddleware, middlewares.SkipAuthorization(func(r *http.Request) bool {
			return r.Method == http.MethodGet
		}),
	)

	r.Use(authorizationMiddleware)

	config.InitRoutes(r, db)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}
	r.Run(fmt.Sprintf(":%d", port))
}
