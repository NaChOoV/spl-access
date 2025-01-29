package controller

import (
	"spl-access/src/dto"
	"spl-access/src/helpers"
	"spl-access/src/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AccessController struct {
	Validator     *validator.Validate
	AccessService *service.AccessService
}

func NewAccessController(accessService *service.AccessService) *AccessController {
	return &AccessController{
		// TODO: Move validator to a global DI scope
		Validator:     validator.New(),
		AccessService: accessService,
	}
}

func (a *AccessController) UpdateOrCreateAccess(c *fiber.Ctx) error {

	// Validate request body
	var accessRequest dto.AccessArrayDto
	if err := c.BodyParser(&accessRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	// Validate the request struct
	if err := a.Validator.Struct(accessRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := a.AccessService.UpdateOrCreateAccess(accessRequest)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a *AccessController) GetAccess(c *fiber.Ctx) error {
	access := a.AccessService.GetAccess()

	return c.JSON(fiber.Map{"data": access})
}
