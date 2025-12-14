package store

import (
	"context"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"uniqueIndex" json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *gorm.DB
}


func NewRoleStore(db *gorm.DB) *RoleStore {
	return &RoleStore{db: db}
}


func (s *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := s.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&role).
		Error

		if err != nil {
		return nil, err
	}

	return &role, nil
}
