package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pangdfg/gopher-social/internal/store"
)

type Storage struct {
	UserCache interface {
		Get(context.Context, uint) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, uint)
	}	
	PostCache interface {
		Get(context.Context, uint) (*store.Post, error)
		Set(context.Context, *store.Post) error
		Delete(context.Context, uint)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		UserCache: &UserStore{rdb: rbd},
		PostCache: &PostStore{rdb: rbd},
	}
}

