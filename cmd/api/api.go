package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/pangdfg/gopher-social/doc"
	"github.com/pangdfg/gopher-social/internal/env"
)


func mount(c *fiber.App, app *application) {

	//Middlewares
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

	//API v1 routes
	v1 := c.Group("/v1")

	//Ops routes
	v1.Get("/health", app.healthCheckHandler)

	//Swagger docs
	v1.Get("/swagger/*", swagger.New())

	//feed
	v1.Get("/feed", app.getUserFeedHandler)

	//Users routes
	users := v1.Group("/users")
	
	users.Put("/activate/:token", app.activateUserHandler)

	users.Use(app.AuthTokenMiddleware)
	users.Put("/update-username", app.updateUsernameHandler)
	users.Put("/change-password", app.ChangePasswordHandler)
	
	user := users.Group("/:userID")

	user.Get("/", app.getUserHandler)
	//user.Put("/follow", app.followUserHandler)
	//user.Put("/unfollow", app.unfollowUserHandler)
	

	//Auth routes
	auth := v1.Group("/auth")
	auth.Post("/user", app.registerUserHandler)
	auth.Post("/token", app.createTokenHandler)

	//Posts routes
	posts := v1.Group("/posts", app.AuthTokenMiddleware)

	posts.Post("/", app.createPostHandler)

	post := posts.Group("/:postID", app.postsContextMiddleware)

	post.Get("/", app.getPostHandler)
	
	post.Post("/", app.createCommentHandler)
	post.Patch("/",  app.updatePostHandler)
	post.Delete("/", app.deletePostHandler)
}