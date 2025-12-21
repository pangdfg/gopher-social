package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pangdfg/gopher-social/internal/store"
)
	
type PaginatedFeedQuery struct {
		Limit  int
		Offset int
		Sort   string
		Tags   []string
		Search string
}

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router		    /feed [get]
func (app *application) getFeedHandler(c *fiber.Ctx) error {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}

	var err error
	fq, err = fq.Parse(c) 
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(fq); err != nil {
		return app.badRequestResponse(c, err)
	}

	feed, err := app.store.Posts.GetFeed(c.Context(), fq)
	if err != nil {
		return app.internalServerError(c, err)
	}

	response := NewPostListResponse(
		feed,
		fq.Limit,
		fq.Offset,
	)

	return app.jsonResponse(c, fiber.StatusOK, response)
}

func (fq PaginatedFeedQuery) Parse(c *fiber.Ctx) (PaginatedFeedQuery, error) {
	if limit := c.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}

	if offset := c.Query("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}

	if sort := c.Query("sort"); sort != "" {
		fq.Sort = sort
	}

	if search := c.Query("search"); search != "" {
		fq.Search = search
	}

	if tags := c.Query("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	return fq, nil
}
