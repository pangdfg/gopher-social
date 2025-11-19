package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pangdfg/gopher-social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func convertToStoreTags(tags []string) []store.Tag {
	storeTags := make([]store.Tag, len(tags))
	for i, tag := range tags {
		storeTags[i] = store.Tag{Name: tag}
	}
	return storeTags
}

func (app *application) createPostHandler(c *fiber.Ctx) error {
	var payload CreatePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user := c.Locals("user").(*store.User)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    convertToStoreTags(payload.Tags),
		UserID:  user.ID,
	}

	ctx := c.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		return app.internalServerError(c, err)
	}	
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": post,
	})
}

func (app *application) getPostHandler(c *fiber.Ctx) error {
	post := c.Locals("post").(*store.Post)

	comments, err := app.store.Comments.GetByPostID(c.Context(), post.ID)
	if err != nil {
		return app.internalServerError(c, err)
	}

	post.Comments = comments

	return c.Status(fiber.StatusOK).JSON(post)
}

func (app *application) deletePostHandler(c *fiber.Ctx) error {
	idParam := c.Params("postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := app.store.Posts.Delete(c.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(c *fiber.Ctx) error {
	post := c.Locals("post").(*store.Post)

	var payload UpdatePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.updatePost(c, post); err != nil {
		return app.internalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

func (app *application) postsContextMiddleware(c *fiber.Ctx) error {
	idParam := c.Params("postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.internalServerError(c, err)
	}

	post, err := app.store.Posts.GetByID(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	c.Locals("post", post)

	return c.Next()
}


func getPostFromCtx(c *fiber.Ctx) *store.Post {
	post, _ := c.Locals("post").(*store.Post)
	return post
}

func (app *application) updatePost(c *fiber.Ctx, post *store.Post) error {
	if err := app.store.Posts.Update(c.Context(), post); err != nil {
		return err
	}

	app.cacheStorage.Users.Delete(c.Context(), post.UserID)
	return nil
}
