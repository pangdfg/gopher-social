package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pangdfg/gopher-social/internal/store"
)

type UserCacheStore interface {
	Get(context.Context, int64) (*store.User, error)
	Set(context.Context, *store.User) error
	Delete(context.Context, int64)
}

type Storage struct {
	Users UserCacheStore
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}

