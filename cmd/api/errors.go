package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) internalServerError(c *fiber.Ctx, err error) error {
	app.logger.Errorw("internal error",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "the server encountered a problem",
	})
}

func (app *application) forbiddenResponse(c *fiber.Ctx) error {
	app.logger.Warnw("forbidden",
		"method", c.Method(),
		"path", c.Path(),
		"error",
	)

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": "forbidden",
	})
}

func (app *application) badRequestResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnw("bad request",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func (app *application) conflictResponse(c *fiber.Ctx, err error) error {
	app.logger.Errorw("conflict response",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func (app *application) notFoundResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnw("not found error",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "not found",
	})
}

func (app *application) unauthorizedError(c *fiber.Ctx, err error) error {
	app.logger.Warnw("unauthorized error",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "unauthorized",
	})
}

func (app *application) unauthorizedErrorResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnw("unauthorized basic error",
		"method", c.Method(),
		"path", c.Path(),
		"error", err.Error(),
	)

	c.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "unauthorized",
	})
}

func (app *application) rateLimitExceededResponse(c *fiber.Ctx, retryAfter string) error {
	app.logger.Warnw("rate limit exceeded",
		"method", c.Method(),
		"path", c.Path(),
	)

	c.Set("Retry-After", retryAfter)

	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"error": "rate limit exceeded, retry after: " + retryAfter,
	})
}
