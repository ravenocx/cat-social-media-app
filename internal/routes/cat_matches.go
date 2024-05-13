package routes

import (
	"github.com/ravenocx/cat-socialx/internal/controllers"
	"github.com/ravenocx/cat-socialx/internal/middleware"
)

func (i *V1Routes) CatMatchRoutes() {
	route := i.Fiber.Group("/v1/cat/match")

	catMatchController := controllers.New(&controllers.V1Repository{
		Repositories: i.Repositories,
	})

	route.Get("", middleware.JWTProtected(), catMatchController.GetCatMatchRequests)
	route.Post("", middleware.JWTProtected(), catMatchController.CreateCatMatch)
	route.Post("/approve", middleware.JWTProtected(), catMatchController.ApproveCatMatch)
	route.Post("/reject", middleware.JWTProtected(), catMatchController.RejectCatMatch)
	route.Delete("/:id", middleware.JWTProtected(), catMatchController.DeleteCatMatch)

}