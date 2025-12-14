package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pangdfg/gopher-social/internal/store"
)

type UserStore struct {
    rdb *redis.Client
}


func (s *UserStore) Get(ctx context.Context, userID uint) (*store.User, error) {
    cacheKey := fmt.Sprintf("user-%d", userID)
    val, err := s.rdb.Get(ctx, cacheKey).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, store.ErrNotFound
        }
        return nil, err
    }
    var u store.User
    if err := json.Unmarshal([]byte(val), &u); err != nil {
        return nil, err
    }
    return &u, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
    cacheKey := fmt.Sprintf("user-%d", user.ID)
    b, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return s.rdb.Set(ctx, cacheKey, b, time.Minute*15).Err()
}

func (s *UserStore) Delete(ctx context.Context, userID uint) {
    cacheKey := fmt.Sprintf("user-%d", userID)
	s.rdb.Del(ctx, cacheKey)
}