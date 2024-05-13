package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ravenocx/cat-socialx/internal/models"
	"github.com/ravenocx/cat-socialx/internal/utils"
)

func (i *V1Repository) AddNewCat(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	expires := claims.Expires

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	newCat := &models.NewCat{}

	if err := c.BodyParser(newCat); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	log.Printf("Payload : %+v", newCat)

	validate := utils.NewValidator()

	if err := validate.Struct(newCat); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	userID := claims.UserID

	cat := &models.Cat{}
	cat.ID = uuid.New()
	cat.UserID = userID
	cat.Name = newCat.Name
	cat.Race = newCat.Race
	cat.Sex = newCat.Sex
	cat.AgeInMonth = newCat.AgeInMonth
	cat.ImageUrls = newCat.ImageUrls
	cat.Description = newCat.Description
	cat.HasMatched = false
	cat.CreatedAt = time.Now()

	if err := validate.Struct(cat); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	log.Printf("Data to add for Cat : %+v", cat)

	if err := i.Repositories.CreateCat(cat); err != nil {
		log.Printf("Failed create new cat : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"id":        cat.ID,
			"createdAt": cat.CreatedAt,
		},
	})
}

func (i *V1Repository) GetCats(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	expires := claims.Expires

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	query := ""

	id := c.Query("id")
	race := c.Query("race")
	sex := c.Query("sex")
	hasMatchedStr := c.Query("hasMatched")
	ageInMonthStr := c.Query("ageInMonth")
	ownedStr := c.Query("owned")
	search := c.Query("search")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if id != "" {
		query += fmt.Sprintf(" AND id = '%s'", id)
	}

	if race != "" {
		allowedRaces := map[string]bool{
			"persian":           true,
			"maine coon":        true,
			"siamese":           true,
			"ragdoll":           true,
			"bengal":            true,
			"sphynx":            true,
			"british shorthair": true,
			"abyssinian":        true,
			"scottish fold":     true,
			"birman":            true,
		}
		if allowedRaces[strings.ToLower(race)] {
			query += fmt.Sprintf(" AND race = '%s'", race)
		}
	}

	if sex != "" {
		sex_lower := strings.ToLower(sex)
		if sex_lower == "male" || sex_lower == "female" {
			query += fmt.Sprintf(" AND sex = '%s'", sex_lower)
		}
	}

	if hasMatchedStr == "true" || hasMatchedStr == "false" {
		hasMatched, _ := strconv.ParseBool(hasMatchedStr)
		query += fmt.Sprintf(" AND hasmatched = %t", hasMatched)
	}

	if ageInMonthStr != "" {
		log.Printf("ageInMonth : %+v", ageInMonthStr)
		var comparison = "="
		var value int
		var err error
		if ageInMonthStr[0] == '<' {
			value, err = strconv.Atoi(ageInMonthStr[1:])
			comparison = "<"
		} else if ageInMonthStr[0] == '>' {
			value, err = strconv.Atoi(ageInMonthStr[1:])
			comparison = ">"
		} else if ageInMonthStr[0] == '=' {
			value, err = strconv.Atoi(ageInMonthStr[1:])
			comparison = "="
		}
		if err == nil {
			query += fmt.Sprintf(" AND ageinmonth %s %d", comparison, value)
		}
	}

	if ownedStr == "true" || ownedStr == "false" {
		claims, err := utils.ExtractTokenMetadata(c)
		if err == nil {
			userID := claims.UserID.String()
			query += fmt.Sprintf(" AND user_id = '%s'", userID)
		}
	}

	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE '%s'`, search)
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			query += fmt.Sprintf(" LIMIT %d", limit)
		}
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	}

	log.Printf("Query params : %+v" , query)

	res, err := i.Repositories.GetCatsData(query)
	if err != nil {
		log.Printf("Failed to get cats data : %+v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": err.Error(),
		})
	}

	cats := []models.CatData{}

	cats = append(cats, res...)

	log.Printf("Cats data : %+v", cats)

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    cats,
	})
}

func (i *V1Repository) UpdateCat(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)

	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	expires := claims.Expires

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	catId := c.Params("id")
	catUUID, err := uuid.Parse(catId)

	log.Printf("Id cat to update : %+v", catId)

	if err != nil {
		log.Printf("Error parsing the params : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	foundedCat, err := i.Repositories.GetCatById(catUUID)
	if err != nil || err == sql.ErrNoRows {
		log.Printf("Error find the cat : %+v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat with this ID not found",
		})
	}

	userID := claims.UserID

	if foundedCat[0].UserID == userID {
		cat_update_request := &models.CatUpdateRequest{}

		if err := c.BodyParser(cat_update_request); err != nil {
			log.Printf("Error parsing the payload :%+v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   fiber.ErrBadRequest.Message,
				"message": err.Error(),
			})
		}

		log.Printf("Payload : %+v", cat_update_request)
		timeNow := time.Now()
		cat_update_request.UpdatedAt = &timeNow

		validate := utils.NewValidator()

		if err := validate.Struct(cat_update_request); err != nil {
			log.Printf("Payload doesn't pass validation : %+v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   fiber.ErrBadRequest.Message,
				"message": utils.ValidatorErrors(err),
			})
		}

		if err := i.Repositories.UpdateCat(foundedCat[0].ID, cat_update_request); err != nil {
			log.Printf("Failed update cat : %+v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fiber.ErrInternalServerError.Message,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "successfully updated cat",
		})
	} else {
		log.Println("Permission denied, only owner can update te cat")
		log.Printf("Id cat to update : %+v", catId)
		log.Printf("Id of Owner cat to update : %+v", foundedCat[0].UserID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   fiber.ErrForbidden.Message,
			"message": "permission denied, only owner can update this cat",
		})
	}
}

func (i *V1Repository) DeleteCat(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	expires := claims.Expires

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	userID := claims.UserID.String()

	catID := c.Params("id")

	if err = i.Repositories.DeleteCat(catID, userID); err != nil {
		log.Printf("Failed to delete cat data : %+v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Deletion successfull",
	})
}
