package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ravenocx/cat-socialx/internal/models"
	"github.com/ravenocx/cat-socialx/internal/utils"
)

func (i *V1Repository) CreateCatMatch(c *fiber.Ctx) error {
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
	userId := claims.UserID

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	catmatch_request := &models.CatMatchRequest{}

	if err := c.BodyParser(catmatch_request); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	log.Printf("Payload : %+v", catmatch_request)

	validate := utils.NewValidator()

	if err := validate.Struct(catmatch_request); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	issuerCat, err := i.Repositories.GetCatById(catmatch_request.CatIssuerID)
	if err != nil {
		log.Printf("Failed to get cat issuer data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if len(issuerCat) == 0 {
		log.Println("Cat issuer not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat issuer not found, please check your request",
		})
	}

	if issuerCat[0].UserID != userId {
		log.Println("Cat issuer is not owned by user that logged in")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat issuer needs to be yours",
		})
	}

	matchcat, err := i.Repositories.GetCatById(catmatch_request.CatMatchID)
	if err != nil {
		log.Printf("Failed to get cat match data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if len(matchcat) == 0 {
		log.Println("Cat match not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat match not found, please check your request",
		})
	}

	cat_match, err := i.Repositories.GetCatMatchByCatIds(catmatch_request.CatMatchID, catmatch_request.CatIssuerID)
	if err != nil {
		log.Printf("Failed to get CatMatch data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if len(cat_match) > 0 {
		log.Println("CatMatch already available when the issuer request for match cat")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "this match request already available when the issuer request for your cat, you can response to the request",
		})
	}

	if issuerCat[0].Sex == matchcat[0].Sex {
		log.Println("Both cat is same gender")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "please choose cat with different gender",
		})
	}

	if issuerCat[0].HasMatched || matchcat[0].HasMatched {
		log.Println("One of the cat is already matched")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "one of the cat is already matched",
		})
	}

	if issuerCat[0].UserID == matchcat[0].UserID {
		log.Println("Both cat is owned by same user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "both cat cant be owned by same user",
		})
	}

	catmatch := &models.CatMatch{}

	catmatch.ID = uuid.New()
	catmatch.CatIssuerID = catmatch_request.CatIssuerID
	catmatch.CatMatchID = catmatch_request.CatMatchID
	catmatch.Message = catmatch_request.Message
	catmatch.Status = "pending"

	currentTime := time.Now().Format(time.RFC3339)
	cm_createdAt, err := time.Parse(time.RFC3339, currentTime)
	if err != nil {
		log.Printf("Failed to parse the time : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}


	catmatch.CreatedAt = cm_createdAt

	if err := validate.Struct(catmatch); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	log.Printf("Cat match data to add : %+v", catmatch)

	if err := PublishToRabbitMQ(catmatch); err != nil {
		log.Printf("Failed to publish to rabbitmq : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if err := i.Repositories.CreateCatMatch(catmatch); err != nil {
		log.Printf("Failed create new CatMatch : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    catmatch,
	})
}

func (i *V1Repository) GetCatMatchRequests(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		log.Printf("Failed to extact the token : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	userId := claims.UserID
	expires := claims.Expires

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	cats, err := i.Repositories.GetCatsByUserId(userId)
	if err != nil {
		log.Printf("Failed to get cats data by userID : %+v", err )
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	log.Printf("Cat that owned by %+v : %+v" , userId, cats)

	userCatMatches := []models.CatMatchDetail{}

	for _, cat := range cats {
		catMatches, err := i.Repositories.GetCatMatchRequests(cat.ID)
		if err != nil {
			log.Printf("Error get CatMatch : %+v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fiber.ErrInternalServerError.Message,
				"message": err.Error(),
			})
		}

		userCatMatches = append(userCatMatches, catMatches...)
	}

	log.Printf("CatMatches : %v", userCatMatches)

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    userCatMatches,
	})
}

func (i *V1Repository) ApproveCatMatch(c *fiber.Ctx) error {
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
	userId := claims.UserID

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	updateRequest := &models.CatMatchUpdateRequest{}

	if err := c.BodyParser(updateRequest); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(updateRequest); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	log.Printf("Payload : %+v", updateRequest)

	cat_match, err := i.Repositories.GetCatMatchById(updateRequest.ID)
	if err != nil {
		log.Printf("Failed to get CatMatch : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if len(cat_match) == 0 {
		log.Println("CatMatch not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat match not found",
		})
	}

	matchCat, err := i.Repositories.GetCatById(cat_match[0].CatMatchID)
	if err != nil {
		log.Printf("Failed to get match cat data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if matchCat[0].UserID != userId {
		log.Println("This request match cat is not owned by the user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "match cat needs to be yours",
		})
	}

	if cat_match[0].Status != "pending" {
		log.Println("This request is not pending")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "this request already approved/rejected",
		})
	}

	issuerCat, err := i.Repositories.GetCatById(cat_match[0].CatIssuerID)
	if err != nil {
		log.Printf("Failed to get issuer cat data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if matchCat[0].HasMatched || issuerCat[0].HasMatched {
		log.Println("One of the cat is alreade matched")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "either match cat or issuer cat already matched with other cat",
		})
	}

	if err := i.Repositories.UpdateCatMatch(updateRequest.ID, "approved"); err != nil {
		log.Printf("Failed to update CatMatch status : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if err := i.Repositories.UpdateCatHasMatched(issuerCat[0].ID); err != nil {
		log.Printf("Failed to update Issuer Cat HasMatched : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if err := i.Repositories.UpdateCatHasMatched(matchCat[0].ID); err != nil {
		log.Printf("Failed to update Match Cat HasMatched : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	issuerCat_Matches, err := i.Repositories.GetCatMatchRequests(issuerCat[0].ID)
	if err != nil {
		log.Printf("Failed to Get Cat Match for Issuer cat : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	for _, issuerCat_match := range issuerCat_Matches {
		if err := i.Repositories.DeleteCatMatchExceptNotPending(issuerCat_match.ID); err != nil {
			log.Printf("Failed to Delete CatMatch that related to Issuer cat : %+v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fiber.ErrInternalServerError.Message,
				"message": err.Error(),
			})
		}
	}

	matchCat_matches, err := i.Repositories.GetCatMatchRequests(matchCat[0].ID)
	if err != nil {
		log.Printf("Failed to Get Cat Match for Match cat : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	for _, matchCat_match := range matchCat_matches {
		if err := i.Repositories.DeleteCatMatchExceptNotPending(matchCat_match.ID); err != nil {
			log.Printf("Failed to Delete CatMatch that related to Match cat : %+v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fiber.ErrInternalServerError.Message,
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success accepted cat match",
	})
}

func (i *V1Repository) RejectCatMatch(c *fiber.Ctx) error {
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
	userId := claims.UserID

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}

	updateRequest := &models.CatMatchUpdateRequest{}

	if err := c.BodyParser(updateRequest); err != nil {
		log.Printf("Error parsing the payload :%+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": err.Error(),
		})
	}

	validate := utils.NewValidator()

	if err := validate.Struct(updateRequest); err != nil {
		log.Printf("Payload doesn't pass validation : %+v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": utils.ValidatorErrors(err),
		})
	}

	log.Printf("Payload : %+v", updateRequest)

	cat_match, err := i.Repositories.GetCatMatchById(updateRequest.ID)
	if err != nil {
		log.Printf("Failed to get CatMatch : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if len(cat_match) == 0 {
		log.Println("CatMatch not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat match not found",
		})
	}

	matchCat, err := i.Repositories.GetCatById(cat_match[0].CatMatchID)
	if err != nil {
		log.Printf("Failed to get match cat data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if matchCat[0].UserID != userId {
		log.Println("This request match cat is not owned by the user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "match cat needs to be yours",
		})
	}

	if cat_match[0].Status != "pending" {
		log.Println("This request is not pending")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "this request already approved/rejected",
		})
	}

	issuerCat, err := i.Repositories.GetCatById(cat_match[0].CatIssuerID)
	if err != nil {
		log.Printf("Failed to get issuer cat data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if matchCat[0].HasMatched || issuerCat[0].HasMatched {
		log.Println("One of the cat is alreade matched")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "either match cat or issuer cat already matched with other cat",
		})
	}

	if err := i.Repositories.UpdateCatMatch(updateRequest.ID, "rejected"); err != nil {
		log.Printf("Failed to update CatMatch status : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success rejected cat match",
	})
}

func (i *V1Repository) DeleteCatMatch(c *fiber.Ctx) error {
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
	userId := claims.UserID

	if now > expires {
		log.Printf("Token already expired, please renew the token : %+v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   fiber.ErrUnauthorized.Message,
			"message": err.Error(),
		})
	}


	id := c.Params("id")
	catMatchId, err := uuid.Parse(id)

	log.Printf("CatMatch id to delete : %+v", id)
	if err != nil {
		log.Printf("Failed to parse the catmatch id params : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	cat_match, err := i.Repositories.GetCatMatchById(catMatchId)
	if err != nil {
		log.Printf("Failed to get CatMatch data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}
	log.Println("test")

	if len(cat_match) == 0 {
		log.Println("CatMatch data not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   fiber.ErrNotFound.Message,
			"message": "cat match not found",
		})
	}

	issuerCat, err := i.Repositories.GetCatById(cat_match[0].CatIssuerID)
	if err != nil {
		log.Printf("Failed to get Cat Issuer data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	if issuerCat[0].UserID != userId {
		log.Println("Issuer cat in this CatMatch data is not owned by the user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "you cant delete other user's match cat request",
		})
	}

	if cat_match[0].Status != "pending" {
		log.Println("This CatMatch request status is not pending")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fiber.ErrBadRequest.Message,
			"message": "this request already approved/rejected and cant be deleted",
		})
	}

	if err := i.Repositories.DeleteCatMatchById(catMatchId); err != nil {
		log.Printf("Failed to delete CatMatch data : %+v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError.Message,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "success deleted cat match",
	})
}