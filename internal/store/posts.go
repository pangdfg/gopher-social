package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        int64           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255" json:"title"`
	Content   string         `json:"content"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Tags      []Tag          `gorm:"many2many:post_tags;" json:"tags"`
	Comments  []Comment      `json:"comments"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:50" json:"name"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *gorm.DB
}


func NewPostStore(db *gorm.DB) *PostStore {
	return &PostStore{db: db}
}


func (s *PostStore) Create(ctx context.Context, post *Post) error {
	return s.db.Create(post).Error
}


func (s *PostStore) GetByID(ctx context.Context, id uint) (*Post, error) {
	post := &Post{}
	err := s.db.Preload("User").Preload("Tags").Preload("Comments").First(post, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}


func (s *PostStore) Update(ctx context.Context,post *Post) error {
	tx := s.db.Model(&Post{}).
		Where("id = ? AND version = ?", post.ID, post.Version).
		Updates(map[string]interface{}{
			"title":   post.Title,
			"content": post.Content,
			"version": post.Version + 1,
		}).First(post)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return tx.Error
	}

	post.Version += 1
	return nil
}

// Delete a post
func (s *PostStore) Delete(ctx context.Context, postID uint) error {
	tx := s.db.Delete(&Post{}, postID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID uint, search string, tags []string, limit, offset int, sort string) ([]Post, error) {
	var posts []Post
	query := s.db.Preload("User").Preload("Tags").Preload("Comments")

	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if len(tags) > 0 {
		query = query.Joins("JOIN post_tags pt ON pt.post_id = posts.id").
			Joins("JOIN tags t ON t.id = pt.tag_id").
			Where("t.name IN ?", tags)
	}

	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	err := query.Order("created_at " + sort).Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}
