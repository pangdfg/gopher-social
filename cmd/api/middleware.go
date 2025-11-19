package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/pangdfg/gopher-social/internal/store"
)

func (app *application) AuthTokenMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return app.unauthorizedError(c, fmt.Errorf("authorization header is missing"))
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return app.unauthorizedError(c, fmt.Errorf("authorization header is malformed"))
	}

	token := parts[1]
	jwtToken, err := app.authenticator.ValidateToken(token)
	if err != nil {
		return app.unauthorizedError(c, err)
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
	if err != nil {
		return app.unauthorizedError(c, err)
	}

	user, err := app.getUser(c.Context(), uint(userID))
	if err != nil {
		return app.unauthorizedError(c, err)
	}

	c.Locals("user", user)
	return c.Next()
}


func (app *application) BasicAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return app.unauthorizedError(c, fmt.Errorf("authorization header is missing"))
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return app.unauthorizedError(c, fmt.Errorf("authorization header malformed"))
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return app.unauthorizedError(c, err)
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return app.unauthorizedError(c, fmt.Errorf("invalid format"))
	}

	username := app.config.auth.basic.user
	pass := app.config.auth.basic.pass

	if creds[0] != username || creds[1] != pass {
		return app.unauthorizedError(c, fmt.Errorf("invalid credentials"))
	}

	return c.Next()
}

func (app *application) checkRolePrecedence(c *fiber.Ctx, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(c.Context(), roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level <= role.Level, nil
}

func (app *application) CheckPostOwnership( requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*store.User)
		post := c.Locals("post").(*store.Post)

		if post.UserID == user.ID {
			return c.Next()
		}

		allowed, err := app.checkRolePrecedence(c , user, requiredRole)
		if err != nil {
			return app.internalServerError(c, err)
		}

		if !allowed {
			return app.forbiddenResponse(c)
		}

		return c.Next()
	}
}

func (app *application) RateLimiterMiddleware(c *fiber.Ctx) error {
	if !app.config.rateLimiter.Enabled {
		return c.Next()
	}

	allow, retryAfter := app.rateLimiter.Allow(c.IP())
	if !allow {
		return app.rateLimitExceededResponse(c, retryAfter.String())
	}

	return c.Next()
}

func (app *application) getUser(ctx context.Context, id uint) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByID(ctx, id)
	}
	user, err := app.cacheStorage.Users.Get(ctx, id)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}

		user, err = app.store.Users.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		err = app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}