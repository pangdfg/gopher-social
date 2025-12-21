package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"size:255" json:"title"`
	Posts     []Post 		 `gorm:"many2many:post_tags;" json:"posts"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TagStore struct {
	db *gorm.DB
}

func NewTagStore(db *gorm.DB) *TagStore {
	return &TagStore{db: db}
}

func (s *TagStore) Create(ctx context.Context, tag *Tag) error {
	var existing Tag
	err := s.db.WithContext(ctx).
		Where("title = ?", tag.Title).
		First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.db.WithContext(ctx).Create(tag).Error
		}
		return err
	}
	tag.ID = existing.ID

	return nil
}

func (s *TagStore) GetByID(ctx context.Context, id uint) (*Tag, error) {
	tag := &Tag{}
	err := s.db.WithContext(ctx).
				First(tag, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return tag, nil
}

func (s *TagStore) Get(ctx context.Context, fq PaginatedFeedQuery) ([]Tag , error) {

	var tags []Tag
	
	query := s.db.WithContext(ctx)

	if fq.Search != "" {
		query = query.Where("title ILIKE ?", "%"+fq.Search+"%")
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

	query = query.Order("tags.created_at " + fq.Sort)

	err := query.
		Limit(fq.Limit).
		Offset(fq.Offset).
		Find(&tags).Error

	if err != nil {
		return nil, err
	}

	return tags, nil
}


func (s *TagStore) Delete(ctx context.Context, tagId uint) error {
	tx := s.db.WithContext(ctx).Delete(&Tag{}, tagId)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}