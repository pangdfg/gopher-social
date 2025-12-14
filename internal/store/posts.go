package store

import (
	"context"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"size:255" json:"title"`
	Content   string         `json:"content"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Tags      pq.StringArray `gorm:"type:varchar(100)[]" json:"tags"`
	Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
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
	return s.db.WithContext(ctx).Create(post).Error
}


func (s *PostStore) GetByID(ctx context.Context, id uint) (*Post, error) {
	post := &Post{}
	err := s.db.WithContext(ctx).
				Preload("User").
				Preload("Comments").
				Preload("Comments.User").
				First(post, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}


func (s *PostStore) Update(ctx context.Context,post *Post) error {
	tx := s.db.WithContext(ctx).Model(&Post{}).
		Where("id = ? AND version = ?", post.ID, post.Version).
		Updates(map[string]interface{}{
			"title":   post.Title,
			"content": post.Content,
			"version": post.Version + 1,
		})

	if tx.Error != nil {
		return tx.Error
	}

	post.Version++
	return nil
}

// Delete a post
func (s *PostStore) Delete(ctx context.Context, postID uint) error {
	tx := s.db.WithContext(ctx).Delete(&Post{}, postID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, fq PaginatedFeedQuery) ([]Post, error) {
	var posts []Post
	
	query := s.db.WithContext(ctx).
		Preload("User").
		Preload("Comments").
		Preload("Comments.User")

	if fq.Search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+fq.Search+"%", "%"+fq.Search+"%")
	}

	if len(fq.Tags) > 0 {
		query = query.Joins("JOIN post_tags pt ON pt.post_id = posts.id").
			Joins("JOIN tags t ON t.id = pt.tag_id").
			Where("t.name IN ?", fq.Tags).
			Distinct("posts.id") // Ensure we don't get duplicate posts if multiple tags match
	}

	if fq.Since != "" {
		query = query.Where("created_at >= ?", fq.Since)
	}
	if fq.Until != "" {
		query = query.Where("created_at <= ?", fq.Until)
	}

	if fq.Sort == "" {
		fq.Sort = "desc"
	}

	query = query.Order("posts.created_at " + fq.Sort)

	err := query.
		Limit(fq.Limit).
		Offset(fq.Offset).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}