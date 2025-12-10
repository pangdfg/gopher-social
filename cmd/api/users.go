package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/pangdfg/gopher-social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil || userID < 1 {
		return app.badRequestResponse(c, err)
	}

	user, err := app.getUser(c.Context(), uint(userID))
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": user,
	})
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(c *fiber.Ctx) error {
	followerUser := getUserFromContext(c)

	followedID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
	}


	err = app.store.Followers.Follow(c.Context(), uint(followedID), followerUser.ID)
	if err != nil {
		switch err {
		case store.ErrConflict:
			return app.conflictResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(c *fiber.Ctx) error {
	unfollowerUser := getUserFromContext(c)

	followedID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
	}


	err = app.store.Followers.Unfollow(c.Context(), uint(followedID), unfollowerUser.ID)
	if err != nil {
		switch err {
		case store.ErrConflict:
			return app.conflictResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(c *fiber.Ctx) error {
	token := c.Params("token")

	err := app.store.Users.Activate(c.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}



func getUserFromContext(c *fiber.Ctx) *store.User {
	user, _ := c.Locals(userCtx).(*store.User)
	return user
}
