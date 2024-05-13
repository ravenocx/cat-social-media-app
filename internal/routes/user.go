package routes

import (
	"github.com/ravenocx/cat-socialx/internal/controllers"
	"github.com/ravenocx/cat-socialx/internal/middleware"
)

func (i *V1Routes) UserRoutes() {
	
	route := i.Fiber.Group("/v1")

	userController := controllers.New(&controllers.V1Repository{
		Repositories: i.Repositories,
	})

	route.Post("/user/register", userController.UserSignUp)
	route.Post("/user/login", userController.UserSignIn)
	route.Post("/token/renew", middleware.JWTProtected(), userController.RenewTokens)

}