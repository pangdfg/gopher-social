package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		GetByID(ctx context.Context, id uint) (*Post, error)
		Create(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id uint) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, fq PaginatedFeedQuery) ([]Post, error)
		GetOneUserFeed(ctx context.Context, fq PaginatedFeedQuery, UserID uint) ([]Post, error)
	}
	Users interface {
		GetByID(ctx context.Context, id uint) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		Create(ctx context.Context, u *User, plain string) error
		Activate(ctx context.Context, userID uint) error
		UpdateUsername(ctx context.Context, user *User) error
		UpdatePassword(ctx context.Context, user *User, plain string) error
		Delete(ctx context.Context, id uint) error
	}
	Comments interface {
		Create(ctx context.Context, c *Comment) error
		GetByPostID(ctx context.Context, postID uint) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, userID, followerID uint) error
		Unfollow(ctx context.Context, followerID, userID uint) error
	}
	Roles interface {
		GetByName(ctx context.Context, name string) (*Role, error)
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
		Roles:     &RoleStore{db: db},
	}
}
