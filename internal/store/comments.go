package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
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

// GetByPostID fetches comments for a given post
func (s *CommentStore) GetByPostID(ctx context.Context, postID uint) ([]Comment, error) {
	var comments []Comment
	err := s.db.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Create a new comment
func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	return s.db.Create(comment).Error
}