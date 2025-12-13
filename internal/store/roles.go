package store

import (
	"context"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
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
	role := &Role{}
	if err := s.db.Where("name = ?", name).First(role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil 
		}
		return nil, err
	}
	return role, nil
}
