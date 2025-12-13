package store

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(c *fiber.Ctx) (PaginatedFeedQuery, error) {

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			fq.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			fq.Offset = o
		}
	}

	// Sort
	if sort := c.Query("sort"); sort != "" {
		fq.Sort = sort
	}

	// Tags
	if tags := c.Query("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	// Search
	if search := c.Query("search"); search != "" {
		fq.Search = search
	}

	// Since
	if since := c.Query("since"); since != "" {
		fq.Since = parseTime(since)
	}

	// Until
	if until := c.Query("until"); until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC3339)
}