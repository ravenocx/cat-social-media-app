package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/cat-socialx/internal/repositories"
)

type V1Routes struct {
	Fiber        *fiber.App
	Repositories *repositories.DatabaseRepositories
}

type iV1Routes interface {
	CatRoutes()
	UserRoutes()
	CatMatchRoutes()
}

func New(v1Routes *V1Routes) iV1Routes {
	return v1Routes
}
