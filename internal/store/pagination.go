package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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


func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	if limit := qs.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			fq.Limit = l
		}
	}
	if offset := qs.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			fq.Offset = o
		}
	}
	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	}
	if tags := qs.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}
	if search := qs.Get("search"); search != "" {
		fq.Search = search
	}
	if since := qs.Get("since"); since != "" {
		fq.Since = parseTime(since)
	}
	if until := qs.Get("until"); until != "" {
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