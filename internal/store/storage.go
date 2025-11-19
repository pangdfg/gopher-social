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
		GetUserFeed(ctx context.Context, userID uint, search string, tags []string, limit, offset int, sort string) ([]Post, error)
	}
	Users interface {
		GetByID(ctx context.Context, id uint) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		Create(ctx context.Context, u *User) error
		Activate(ctx context.Context, token string) error
		UpdateUsername(ctx context.Context, user *User) error
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

// NewStorage returns a GORM-based Storage
func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
		Roles:     &RoleStore{db: db},
	}
}

// Optional helper for GORM transactions
func withTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}
