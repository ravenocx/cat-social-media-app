package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ravenocx/cat-socialx/internal/models"
	"github.com/ravenocx/cat-socialx/internal/utils"
)

func (i *V1Repository) RenewTokens(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	expiresAccessToken := claims.Expires

	if now > expiresAccessToken {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	renew := &models.Renew{}

	if err := c.BodyParser(renew); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	log.Printf("Payload : %+v", renew)

	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		log.Printf("Failed to parse the refresh token : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	if now < expiresRefreshToken {
		userID := claims.UserID

		user, err := i.Repositories.GetUserByID(userID)
		if err != nil {
			log.Printf("Failed to get user data : %+v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   fiber.ErrNotFound.Message,
				"message": err.Error(),
			})
		}

		tokens, err := utils.GenerateNewTokens(user.ID.String())
		if err != nil {
			log.Printf("Failed to generate new token : %+v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fiber.ErrInternalServerError.Message,
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": "unauthorized, your session was ended earlier",
		})
	}
}
