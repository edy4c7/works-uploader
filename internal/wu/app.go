package wu

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/edy4c7/works-uploader/internal/config"
	"github.com/edy4c7/works-uploader/internal/entities"
	"github.com/edy4c7/works-uploader/internal/i18n"
	"github.com/edy4c7/works-uploader/internal/lib"
	"github.com/edy4c7/works-uploader/internal/middlewares"
	"github.com/gin-gonic/gin"
)

//Run run app
func Run() {
	r := gin.Default()

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(entities.Work{}, entities.Activity{}, entities.User{})

	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	auth0Issuer := os.Getenv("AUTH0_ISSUER")
	auth0JWK := os.Getenv("AUTH0_JWK")

	jwtMiddleware := middlewares.NewJWTMiddleware(auth0Audience, auth0Issuer, auth0JWK)
	authorizationMiddleware := middlewares.NewAuthorizationMiddleware(
		jwtMiddleware, middlewares.SkipAuthorization(func(r *http.Request) bool {
			return r.Method == http.MethodGet
		}),
	)
	r.Use(authorizationMiddleware)

	authenticationMiddleware := middlewares.NewAuthenticationMiddleware(func(r *http.Request) bool {
		if strings.HasSuffix(r.URL.Path, "/users") {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]

			return lib.CheckJWTScope(auth0JWK, "access:users", token)
		}

		return true
	})
	r.Use(authenticationMiddleware)

	r.Use(middlewares.NewErrorMiddleware(i18n.NewPrinter()))

	config.InitRoutes(r, db)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}
	r.Run(fmt.Sprintf(":%d", port))
}
