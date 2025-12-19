package main

import (
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
	Username string`json:"username" validate:"required,max=100"`
	Email	string`json:"email" validate:"required,email,max=255"`
	*store.Role
	Token string `json:"token"`
}


type RegisterUserData struct {
	ID        uint   
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
	IsActive  bool    
	RoleID    uint
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/user [post]
func (app *application) registerUserHandler(c *fiber.Ctx) error {
	var payload RegisterUserPayload
	ctx := c.Context()

	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}
	
	role, err := app.store.Roles.GetByName(ctx, "user")
	if err != nil {
		return app.internalServerError(c, err)
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		IsActive: false,
		RoleID: role.ID,
	}

	claims := jwt.MapClaims{
    "email": payload.Email,
	"id" : uuid.New().String(),
    "exp":   time.Now().Add(24 * time.Hour).Unix(), 
    "type":  "activation",
	}

	token, err := app.authenticator.GenerateToken(claims)

	userWithToken := UserWithToken{
		Username:  user.Username,
		Email: user.Email,
		Role: role,
		Token: token,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, token)
	isProdEnv := app.config.env == "production"

	err = app.store.Users.Create(ctx, user, payload.Password)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(c, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(c, err)
		default:
			app.internalServerError(c, err)
		}
		return err
	}
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

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/token [post]
func (app *application) createTokenHandler(c *fiber.Ctx) error {
	var payload CreateUserTokenPayload

	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user, err := app.store.Users.GetByEmail(c.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}
	if err := user.Authenticate(payload.Password); err != nil {
		return app.unauthorizedErrorResponse(c, err)
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"role": user.Role.Name,
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

type PasswordPayload struct {
	NewPassword string `json:"new_password" validate:"required,min=3,max=72"`
	OldPassword string `json:"old_password" validate:"required,min=3,max=72"`
}

// ChangePasswordHandler godoc
//
//	@Summary		Change user password
//	@Description	Allows an authenticated user to change their password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		Password	true	"Password payload"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error	"Bad request / invalid payload"
//	@Failure		401		{object}	error	"Unauthorized / current password invalid"
//	@Failure		500		{object}	error	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/change-password [put]
func (app *application) ChangePasswordHandler(c *fiber.Ctx) error {
	authUser := getUserFromContext(c)

	var payload PasswordPayload

	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user, err := app.store.Users.GetByID(c.Context(), authUser.ID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}
	
	if err := user.Authenticate(payload.OldPassword); err != nil {
		return app.unauthorizedErrorResponse(c, err)
	}

	if err := app.store.Users.UpdatePassword(c.Context(), authUser, payload.NewPassword); err != nil {
	switch err{
	case store.ErrConflict:
		return app.conflictResponse(c, err)
	default:
		return app.internalServerError(c, err)
		}
	}

	updatedUser, err := app.store.Users.GetByID(c.Context(), authUser.ID)
	if err != nil {
		return app.internalServerError(c, err)
	}


	return c.Status(fiber.StatusOK).JSON(updatedUser)
}