package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(data)
}

func writeJSONError(c *fiber.Ctx, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(c, status, &envelope{Error: message})
}

func (app *application) jsonResponse(c *fiber.Ctx, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(c, status, &envelope{data})
}

func readJSON(c *fiber.Ctx, dst any) error {
	maxBytes := 1_048_576 

	if len(c.Body()) > maxBytes {
		return fiber.ErrRequestEntityTooLarge
	}

	return c.BodyParser(dst)
}
