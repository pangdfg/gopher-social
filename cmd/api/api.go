package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
	fromEmail string
	exp       time.Duration
}

type mailTrapConfig struct {
	apiKey string
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *fiber.App {
	appFiber := fiber.New(fiber.Config{
		AppName:      "Gopher Social API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return app.internalServerError(c, err)
		},
	})

	// ------------ GLOBAL MIDDLEWARE ------------
	appFiber.Use(requestid.New())     // Like chi RequestID
	appFiber.Use(recover.New())       // Panic recover
	appFiber.Use(logger.New())        // Request logger
	appFiber.Use(cors.New(cors.Config{
		AllowOrigins:     env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if app.config.rateLimiter.Enabled {
		appFiber.Use(app.RateLimiterMiddlewareFiber())
	}

	// Timeout middleware (per request)
	appFiber.Use(timeout.New(timeout.Config{
		Timeout: 60 * time.Second,
	}))