package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/cat-socialx/internal/repositories"
)

type V1Repository struct {
	Repositories *repositories.DatabaseRepositories
}

type iV1Controller interface {
	UserSignUp(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
	AddNewCat(c *fiber.Ctx) error
	GetCats(c *fiber.Ctx) error
	UpdateCat(c *fiber.Ctx) error
	DeleteCat(c *fiber.Ctx) error
	CreateCatMatch(c *fiber.Ctx) error
	GetCatMatchRequests(c *fiber.Ctx) error
	ApproveCatMatch(c *fiber.Ctx) error
	RejectCatMatch(c *fiber.Ctx) error
	DeleteCatMatch(c *fiber.Ctx) error
	RenewTokens(c *fiber.Ctx) error
}

func New(v1Repository *V1Repository) iV1Controller {
	return v1Repository
}
