package routes

import (
	"github.com/ravenocx/cat-socialx/internal/controllers"
	"github.com/ravenocx/cat-socialx/internal/middleware"
)

func (i *V1Routes) CatRoutes() {
	route := i.Fiber.Group("/v1/cat")

	catController := controllers.New(&controllers.V1Repository{
		Repositories: i.Repositories,
	})

	route.Get("", middleware.JWTProtected(), catController.GetCats)
	route.Post("", middleware.JWTProtected(), catController.AddNewCat)
	route.Delete("/:id", middleware.JWTProtected(), catController.DeleteCat)
	route.Put("/:id", middleware.JWTProtected(), catController.UpdateCat)
}