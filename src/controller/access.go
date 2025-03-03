package controller

import (
	"spl-access/src/config"
	"spl-access/src/dto"
	"spl-access/src/helpers"
	"spl-access/src/model"
	"spl-access/src/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AccessController struct {
	Validator     *validator.Validate
	accessService *service.AccessService
	config        *config.EnvironmentConfig
}

func NewAccessController(
	accessService *service.AccessService,
	config *config.EnvironmentConfig,
) *AccessController {
	return &AccessController{
		// TODO: Move validator to a global DI scope
		Validator:     validator.New(),
		accessService: accessService,
		config:        config,
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

	err := a.accessService.UpdateOrCreateAccess(accessRequest)
	if err != nil {
		return helpers.InternalError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a *AccessController) GetObfuscateAccess(c *fiber.Ctx) error {
	utcTime := time.Now().UTC()
	var accesses *[]model.Access

	if helpers.IsChileSleepTime(utcTime, a.config.Zone) {
		accesses = &[]model.Access{}
	} else {
		accesses = a.accessService.GetTodayAccess()
	}

	if len(*accesses) == 0 {
		return c.JSON(fiber.Map{"data": []model.Access{}})
	}

	obfuscateAccess := helpers.MaskAccessData(accesses)

	return c.JSON(fiber.Map{"data": obfuscateAccess})
}

func (a *AccessController) GetAccess(c *fiber.Ctx) error {
	access := a.accessService.GetTodayAccess()
	if len(*access) == 0 {
		return c.JSON(fiber.Map{"data": []model.Access{}})
	}

	return c.JSON(fiber.Map{"data": access})
}
