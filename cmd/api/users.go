package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/pangdfg/gopher-social/internal/store"
)

type userKey string

const userCtx userKey = "user"


func (app *application) getUserHandler(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil {
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
