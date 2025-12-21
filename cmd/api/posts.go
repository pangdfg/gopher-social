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

func (app *application) postsContextMiddleware(c *fiber.Ctx) error {
	idParam := c.Params("postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
	}
	var post *store.Post

	if app.config.redisCfg.enabled {
		post, err = app.cacheStorage.PostCache.Get(c.Context(), uint(id))
		if err != nil && err != store.ErrNotFound {
			app.logger.Warnw("cache get failed", "postID", id, "error", err.Error())
			post = nil
		}
	}

	if post == nil {
		post, err = app.store.Posts.GetByID(c.Context(), uint(id))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				return app.notFoundResponse(c, err)
			default:
				return app.internalServerError(c, err)
			}
		}

		c.Locals("post", post)

		if app.config.redisCfg.enabled {
			if err := app.cacheStorage.PostCache.Set(c.Context(), post); err != nil {
				app.logger.Warnw("cache set failed", "postID", post.ID, "error", err.Error())
			}
		}
	} else {
		c.Locals("post", post)
	}

	return c.Next()
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
	if err := readJSON(c, &payload); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid JSON body")
	}

	ctx := c.Context()
	user := c.Locals("user").(*store.User)

	var tags []store.Tag
	for _, t := range payload.Tags {
		tag := store.Tag{Title: t}
		if err := app.store.Tags.Create(ctx, &tag); err != nil {
			return app.internalServerError(c, err)
		}	
		tags = append(tags, tag)
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  user.ID,
		Tags: tags,
	}

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
	postInterface := c.Locals("post")
	post, ok := postInterface.(*store.Post)
	if !ok || post == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "post not found",
		})
	}

	response := NewPostResponse(
		post,
	)
	return c.Status(fiber.StatusOK).JSON(response)
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
	
	response := NewPostResponse(
		updatedPost,
	)

	return app.jsonResponse(c, fiber.StatusOK, response)
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

	if err := readJSON(c, &payload); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid JSON body")
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
	
	return app.jsonResponse(c, fiber.StatusCreated, fiber.Map{
		"data": comment,
	})
}

// getTagHandler godoc
//
// @Summary      Get tag by ID
// @Description  Retrieve a single tag by its ID including related posts.
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tagID   path      int  true  "Tag ID"
// @Success      200     {object}  store.Tag
// @Failure      400     {object}  ErrorResponse  "Invalid tag ID"
// @Failure      404     {object}  ErrorResponse  "Tag not found"
// @Failure      500     {object}  ErrorResponse  "Internal server error"
// @Router       /tags/{tagID} [get]
func (app *application) getTagHandler(c *fiber.Ctx) error {
	idParam := c.Params("tagID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1{
		return app.badRequestResponse(c, err)
	}

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	var tag *store.Tag
	tag, err = app.store.Tags.GetByID(c.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	post, err := app.store.Posts.GetByTagID(c.Context(), fq, tag.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	response := NewTagPostListResponse(
		post,
		tag,
		fq.Limit,
		fq.Offset,
	)

	return app.jsonResponse(c, fiber.StatusOK, response)
}

// getTagTitleHandler godoc
//
// @Summary      List tags
// @Description  Retrieve a paginated list of tags with optional search and sorting.
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        search   query     string  false  "Search by tag title"
// @Param        limit    query     int     false  "Number of items to return"  default(20)
// @Param        offset   query     int     false  "Number of items to skip"    default(0)
// @Param        sort     query     string  false  "Sort order (asc|desc)"       default(desc)
// @Success      200      {array}   store.Tag
// @Failure      400      {object}  ErrorResponse  "Invalid query parameters"
// @Failure      500      {object}  ErrorResponse  "Internal server error"
// @Router       /tags [get]
func (app *application) getTagTitleHandler(c *fiber.Ctx) error {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}

	var err error
	fq, err = fq.Parse(c) 
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	tag, err := app.store.Tags.Get(c.Context(), fq)
	if err != nil {
		return app.internalServerError(c, err)
	}

	response := NewTagsListResponse(
		tag,
		fq.Limit,
		fq.Offset,
	)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteTag godoc
//
//	@Summary		Deletes a tag
//	@Description	Delete a tag by ID
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Tag ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/tags/{id} [delete]
func (app *application) deleteTagHandler(c *fiber.Ctx) error {
	idParam := c.Params("tagID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := app.store.Tags.Delete(c.Context(), uint(id)); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}