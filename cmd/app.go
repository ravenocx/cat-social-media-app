package cmd

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/cat-socialx/config"
	"github.com/ravenocx/cat-socialx/internal/middleware"
	"github.com/ravenocx/cat-socialx/internal/repositories"
	"github.com/ravenocx/cat-socialx/internal/routes"
)

func (i *Http) StartApp() {
	serverConfig := config.FiberConfig()

	app := fiber.New(serverConfig)

	middleware.FiberMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	repo := repositories.New(i.DB)

	route := routes.New(&routes.V1Routes{
		Fiber:        app,
		Repositories: repo,
	})

	route.UserRoutes()
	route.CatRoutes()
	route.CatMatchRoutes()

	if err := app.Listen(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
