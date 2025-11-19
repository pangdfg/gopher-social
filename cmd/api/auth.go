package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pangdfg/gopher-social/internal/mailer"
	"github.com/pangdfg/gopher-social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}


func (app *application) registerUserHandler(c *fiber.Ctx) error {
	var payload RegisterUserPayload

	// Parse JSON
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return app.internalServerError(c, err)
	}

	ctx := c.Context()

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])


	userWithToken := UserWithToken{
		User:  user,
		Token: hashToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"

	mailVars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, mailVars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		return app.internalServerError(c, err)
	}

	app.logger.Infow("Email sent", "status code", status)

	return c.Status(fiber.StatusCreated).JSON(userWithToken)
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createTokenHandler(c *fiber.Ctx) error {
	var payload CreateUserTokenPayload

	// Parse request
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	// Fetch user
	user, err := app.store.Users.GetByEmail(c.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.unauthorizedErrorResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	// Compare password
	if err := user.Password.Compare(payload.Password); err != nil {
		return app.unauthorizedErrorResponse(c, err)
	}

	// JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(token)
}
