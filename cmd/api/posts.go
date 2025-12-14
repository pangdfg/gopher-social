package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pangdfg/gopher-social/internal/store"
)
 
type CreateCommentPayload struct {
	PostID  uint   `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,max=500"`
}

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(c *fiber.Ctx) error {
	var payload CreatePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user := c.Locals("user").(*store.User)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := c.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		return app.internalServerError(c, err)
	}	
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": payload,
	})
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(c *fiber.Ctx) error {
	post := c.Locals("post").(*store.Post)

	comments, err := app.store.Comments.GetByPostID(c.Context(), post.ID)
	if err != nil {
		return app.internalServerError(c, err)
	}

	post.Comments = comments

	return c.Status(fiber.StatusOK).JSON(post)
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(c *fiber.Ctx) error {
	post, ok := c.Locals("post").(*store.Post)
	if !ok || post == nil {
		return app.internalServerError(c, errors.New("post context missing"))
	}

	authUser := getUserFromContext(c)
	isOwner := post.UserID == authUser.ID
	if !isOwner {
		return app.forbiddenResponse(c)
	}

	if err := app.store.Posts.Delete(c.Context(), uint(post.ID)); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(c *fiber.Ctx) error {

	post, ok := c.Locals("post").(*store.Post)
	if !ok || post == nil {
		return app.internalServerError(c, errors.New("post context missing"))
	}

	authUser := getUserFromContext(c)
	isOwner := post.UserID == authUser.ID
	if !isOwner {
		return app.forbiddenResponse(c)
	}

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

	if err := app.store.Posts.Update(c.Context(), post); err != nil {
	switch {
	case errors.Is(err, store.ErrConflict):
		return app.conflictResponse(c, err)
	default:
		return app.internalServerError(c, err)
	}
	}
	updatedPost, err := app.store.Posts.GetByID(c.Context(), post.ID)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedPost)
}

func (app *application) postsContextMiddleware(c *fiber.Ctx) error {
	idParam := c.Params("postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
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

// CreateComment godoc
//
//	@Summary		Creates a comment
//	@Description	Creates a comment
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCommentPayload	true	"Post payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [post]
func (app *application) createCommentHandler(c *fiber.Ctx) error {
	var payload CreateCommentPayload
	if err := c.BodyParser(&payload); err != nil {
		return app.badRequestResponse(c, err)
	}
	
	idParam := c.Params("postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.internalServerError(c, err)
	}

	user := c.Locals("user").(*store.User)
	comment := &store.Comment{
		PostID:  uint(id),
		Content: payload.Content,
		UserID:  user.ID,
		User: *user,
	}

	ctx := c.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		return app.internalServerError(c, err)
	}	
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": comment,
	})
}
