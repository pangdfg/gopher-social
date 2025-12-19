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
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Username  string    `gorm:"uniqueIndex;size:255" json:"username"`
	Password []byte `gorm:"column:password;not null" json:"-"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
	RoleID    uint 
	Role   	  Role    `gorm:"foreignKey:RoleID;references:ID"`
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

func (u *User) Authenticate(plain string) error {
	p := password{
		hash: []byte(u.Password),
	}
	return p.Compare(plain)
}



func (s *UserStore) Create(ctx context.Context,user *User, plain string) error {
	var p password
	if err := p.Set(plain); err != nil {
		return err
	}

	user.Password = p.hash

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}


func (s *UserStore) GetByID(ctx context.Context, userID uint) (*User, error) {
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
	if err := s.db.Preload("Role").Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}


func (s *UserStore) Activate(ctx context.Context, userID uint) error {
	return s.db.Model(&User{}).Where("id = ? ", userID).Update("is_active", true).Error
}

func (s *UserStore) Delete(ctx context.Context, userID uint) error {
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
	tx := s.db.Model(&User{}).Where("id = ? AND is_active = ?", user.ID, true).Updates(map[string]interface{}{
		"Username": user.Username,
	})
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

func (s *UserStore) UpdatePassword(ctx context.Context, user *User, plain string) error {
	var p password
	if err := p.Set(plain); err != nil {
		return err
	}

	tx := s.db.Model(&User{}).Where("id = ? AND is_active = ?", user.ID, true).Updates(map[string]interface{}{
		"Password": p.hash,
	})

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return tx.Error
	}
	return nil
}