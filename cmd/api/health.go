package main

import (
	"github.com/gofiber/fiber/v2"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	data := fiber.Map{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	return c.Status(fiber.StatusOK).JSON(data)
}