package store

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
	ErrNotFound          = errors.New("not found")
)

// User model
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Username  string    `gorm:"uniqueIndex;size:255" json:"username"`
	Password  password  `gorm:"-" json:"-"`
	PasswordHash string `gorm:"size:255" json:"-"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
	RoleID    int64
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}


func (s *UserStore) Create(ctx context.Context,user *User) error {
	if err := user.Password.Set(*user.Password.text); err != nil {
		return err
	}
		
	user.PasswordHash = string(user.Password.hash)
	if err := s.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}


func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	user := &User{}
	err := s.db.WithContext(ctx).
		Preload("Role").
		Where("id = ? AND is_active = ?", userID, true).
		First(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}


func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	if err := s.db.Preload("Role").Where("email = ? AND is_active = ?", email, true).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}


func (s *UserStore) Activate(ctx context.Context, userID string) error {
	return s.db.Model(&User{}).Where("id = ?", userID).Update("is_active", true).Error
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	tx := s.db.Delete(&User{}, userID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *UserStore) UpdateUsername(ctx context.Context, user *User) error {
	tx := s.db.Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"Username": user.Username,
	}).First(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateUsername
		}
		return tx.Error
	}
	return nil
}