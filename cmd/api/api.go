package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/pangdfg/gopher-social/internal/env"
)

func mount(c *fiber.App, app *application) {
	c.Use(requestid.New())     
	c.Use(recover.New())      
	c.Use(logger.New())      
	c.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if app.config.rateLimiter.Enabled {
		c.Use(app.RateLimiterMiddleware)
	}

	v1 := c.Group("/v1")

	v1.Get("/health", app.healthCheckHandler)
	
	users := v1.Group("/users")

	users.Put("/activate/{token}", app.activateUserHandler)

	users.Use(app.AuthTokenMiddleware)

	user := users.Group(("/:userID"))

	user.Get("/", app.getUserHandler)
	user.Put("/follow", app.followUserHandler)
	user.Put("/unfollow", app.unfollowUserHandler)

	//users.Get("/feed", app.getUserFeedHandler)

	auth := v1.Group("/auth")
	auth.Post("/user", app.registerUserHandler)
	auth.Post("/token", app.createTokenHandler)
}