package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pangdfg/gopher-social/internal/store"
)

type PostStore struct {
	rdb *redis.Client
}

func (s *PostStore) Get(ctx context.Context, postID uint) (*store.Post, error) {
	cacheKey := fmt.Sprintf("post-%d", postID)
	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	var p store.Post
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostStore) Set(ctx context.Context, post *store.Post) error {
	cacheKey := fmt.Sprintf("post-%d", post.ID)
	b, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, cacheKey, b, time.Minute*2).Err()
}

func (s *PostStore) Delete(ctx context.Context, postID uint) {
	cacheKey := fmt.Sprintf("post-%d", postID)
	s.rdb.Del(ctx, cacheKey)
}

func NewPostStore(rdb *redis.Client) *PostStore {
	return &PostStore{rdb: rdb}
}
