package dsw

import "github.com/gin-gonic/gin"

// import (
// 	"context"

// 	firebase "firebase.google.com/go"
// 	"github.com/gin-gonic/gin"
// 	"github.com/edy4c7/darkpot-school-works/internal/middlewares"
// 	"github.com/edy4c7/darkpot-school-works/internal/infrastructures/baas"
// )

//Run run app
func Run() {
	r := gin.Default()

	// cc := context.Background()
	// fbApp, err := firebase.NewApp(cc, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// authClient, err := fbApp.Auth(cc)
	// if err != nil {
	// 	panic(err)
	// }

	// authMiddleware := middlewares.NewAuthenticationMiddleware(
	// 	infrastructures.NewAuthServiceImpl(authClient))

	r.Run()
}
