package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	f := fiber.New()

	// Middlewares
	f.Use(recover.New())
	f.Use(logger.New())

	// Routes
	v1 := f.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Gopher Social API v1")
	})
	/*
	// Auth
	v1.Post("/auth/login", app.handleLogin)
	v1.Post("/auth/register", app.handleRegister)

	// Posts
	v1.Get("/posts", app.handleListPosts)
	v1.Post("/posts", app.handleCreatePost)

	// Users
	v1.Get("/users/:username", app.handleGetUser)
*/
	log.Fatal(f.Listen(":3000"))
}
