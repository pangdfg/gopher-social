package store

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrConflict = errors.New("already following")

type Follower struct {
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	FollowerID uint      `gorm:"primaryKey" json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerStore struct {
	db *gorm.DB
}

func NewFollowerStore(db *gorm.DB) *FollowerStore {
	return &FollowerStore{db: db}
}

// Follow creates a new follower record
func (s *FollowerStore) Follow(ctx context.Context, followerID uint, userID uint) error {
	f := Follower{
		UserID:     userID,
		FollowerID: followerID,
		CreatedAt:  time.Now(),
	}

	err := s.db.WithContext(ctx).Create(&f).Error
	if err != nil {
		// handle unique constraint conflict (duplicate primary key)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrConflict
		}
		return err
	}

	return nil
}

// Unfollow deletes a follower record
func (s *FollowerStore) Unfollow(ctx context.Context, followerID uint, userID uint) error {
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND follower_id = ?", userID, followerID).
		Delete(&Follower{}).Error
	return err
}