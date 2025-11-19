package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/pangdfg/gopher-social/internal/env"
)

func mount(app *fiber.App, c *application) {
	app.Use(requestid.New())     
	app.Use(recover.New())      
	app.Use(logger.New())      
	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if c.config.rateLimiter.Enabled {
		app.Use(c.RateLimiterMiddleware)
	}

	v1 := app.Group("/v1")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Gopher Social API v1")
	})

	v1.Get("/health", c.healthCheckHandler)
}