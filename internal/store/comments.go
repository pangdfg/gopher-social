package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"-"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}


type CommentStore struct {
	db *gorm.DB
}

// NewCommentStore creates a new CommentStore
func NewCommentStore(db *gorm.DB) *CommentStore {
	return &CommentStore{db: db}
}


func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
    return s.db.WithContext(ctx).Create(comment).Error
}
	