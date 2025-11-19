package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	data := fiber.Map{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	return c.Status(fiber.StatusOK).JSON(data)
}