package main

import (
	"time"

	"github.com/pangdfg/gopher-social/internal/store"
)

type PostResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        UserMini  `json:"author"`
	Tags          []TagsMini`json:"tags"`
	Comments      []CommentMini `json:"commnts"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserMini struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type TagsMini struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type FeedResponse struct {
	Posts []PostMini     `json:"posts"`
	PostsCount int        `json:"posts_count"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type PostMini struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        UserMini  `json:"author"`
	Tags          []TagsMini`json:"tags"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type TagsListResponse struct {
	Tags   []TagsMini `json:"tags"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type TagsResponse struct {
	Tag   TagsMini     `json:"tag"`
	Posts FeedResponse `json:"posts"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type CommentMini struct {
	ID uint                 `json:"id"`
	Content       string    `json:"content"`
	Author        UserMini  `json:"author"`
	CreatedAt     time.Time `json:"created_at"`
}
type UserResponse struct {
	ID uint            `json:"id"`
	Email string       `json:"email"`
	Username string    `json:"username"`
	Role string        `json:"role"`
	Posts FeedResponse `json:"posts"`
}

func NewPostListResponse(posts []store.Post, limit, offset int) FeedResponse {
	res := make([]PostMini, 0, len(posts))
	for _, p := range posts {
		tags := make([]TagsMini, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, TagsMini{
				ID:    t.ID,
				Title: t.Title,
			})
		}
		res = append(res, PostMini{
			ID: p.ID,
			Title: p.Title,
			Content: p.Content,
			Author: UserMini{
				ID:       p.User.ID,
				Username: p.User.Username,
				Role:     p.User.Role.Name,
			},
			Tags: tags,
			CommentsCount: len(p.Comments),
			CreatedAt: p.CreatedAt,

		})
	}

	return FeedResponse{
		Posts:  res,
		PostsCount: len(res),
		Limit:  limit,
		Offset: offset,
	}
}

func NewPostResponse(post *store.Post) PostResponse {
	tags := make([]TagsMini, 0, len(post.Tags))
	for _, t := range post.Tags {
		tags = append(tags, TagsMini{
			ID:    t.ID,
			Title: t.Title,
		})
	}

	comments := make([]CommentMini, 0, len(post.Comments))
	for _, c := range post.Comments {
		comments = append(comments, CommentMini{
			ID:    c.ID,
			Content: c.Content,
			Author: UserMini{
				ID:       c.User.ID,
				Username: c.User.Username,
				Role:     c.User.Role.Name,
			},
			CreatedAt: c.CreatedAt,
		})
	}

	return PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Author: UserMini{
			ID:       post.User.ID,
			Username: post.User.Username,
			Role:     post.User.Role.Name,
		},
		Tags:          tags,
		Comments:      comments,
		CommentsCount: len(comments),
		CreatedAt:     post.CreatedAt,
	}
}

func NewUserResponse(user *userWithPosts, limit int, offset int) UserResponse {
	post := make([]PostMini, 0, len(user.Posts))
	for _, p := range user.Posts{
		tags := make([]TagsMini, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, TagsMini{
				ID:    t.ID,
				Title: t.Title,
			})
		}
		post = append(post, PostMini{
			ID: p.ID,
			Title: p.Title,
			Content: p.Content,
			Author: UserMini{
				ID:       p.User.ID,
				Username: p.User.Username,
				Role:     p.User.Role.Name,
			},
			Tags: tags,
			CommentsCount: len(p.Comments),
			CreatedAt: p.CreatedAt,

		})
	}


	return UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role.Name,
		Posts:  FeedResponse{
			Posts:  post,
			PostsCount: len(post),
			Limit:  limit,
			Offset: offset,
	},
	}
}

func NewTagsListResponse(tags []store.Tag, limit int, offset int) TagsListResponse {
	tag := make([]TagsMini, 0, len(tags))
		for _, t := range tags {
			tag = append(tag, TagsMini{
				ID:    t.ID,
				Title: t.Title,
			})
		}

	return  TagsListResponse{
		Tags: tag,
		Limit: limit,
		Offset: offset,
	}
}

func NewTagPostListResponse(posts []store.Post, tag *store.Tag, limit int, offset int) TagsResponse {
	res := NewPostListResponse(posts, limit, offset)
	return TagsResponse{
		Tag: TagsMini{
			ID: tag.ID,
			Title: tag.Title,
		},
		Posts:  res,
		Limit:  limit,
		Offset: offset,
	}
}