package controllers

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ravenocx/cat-socialx/internal/models"
	"github.com/ravenocx/cat-socialx/internal/utils"
)

func (i *V1Repository) UserSignUp(c *fiber.Ctx) error {
	signUp := &models.SignUpRequest{}

	if err := c.BodyParser(signUp); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(signUp); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	log.Printf("Payload : %+v", signUp)

	user := &models.User{}

	user.ID = uuid.New()
	user.Email = signUp.Email
	user.Name = signUp.Name

	now := time.Now().Format(time.RFC3339)
	createdAt, err := time.Parse(time.RFC3339, now)
	if err != nil {
		log.Printf("Failed to parse the time : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	user.CreatedAt = createdAt
	user.UserStatus = 1
	user.UserRole = "user"

	tokens, err := utils.GenerateNewTokens(user.ID.String())
	if err != nil {
		log.Printf("Failed to generate new token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	user.AccessToken = tokens.Access

	// Dont hash the password before validate the struct
	user.Password = signUp.Password

	if err := validate.Struct(user); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}
	user.Password = utils.GeneratePassword(signUp.Password)

	if err := i.Repositories.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			log.Println("Duplicate on email")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   fiber.ErrConflict.Message,
				"message": err.Error(),
			})
		}

		log.Printf("Failed to create new user : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	// Delete password hash field from JSON view.
	user.Password = ""

	responseData := models.AuthResponse{
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: user.AccessToken,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    responseData,
	})
}

func (i *V1Repository) UserSignIn(c *fiber.Ctx) error {
	signIn := &models.SignInRequest{}

	if err := c.BodyParser(signIn); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(signIn); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	log.Printf("Payload : %+v", signIn)

	user, err := i.Repositories.GetUserByEmail(signIn.Email)
	if err != nil {
		log.Printf("Failed to get user data : %+v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": err.Error(),
		})
	}

	if err = utils.ComparePasswords(user.Password, signIn.Password); err != nil {
		log.Printf("Failed to compare password : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	token, err := utils.GenerateNewTokens(user.ID.String())
	if err != nil {
		log.Printf("Failed to generate new token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	responseData := models.AuthResponse{
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: token.Access,
	}

	return c.JSON(fiber.Map{
		"message": "User logged successfully",
		"data":    responseData,
	})
}
