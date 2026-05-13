package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Sula2007/user-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{client: client}
}

func (c *UserCache) Set(ctx context.Context, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "user:"+user.ID, data, 30*time.Minute).Err()
}

func (c *UserCache) Get(ctx context.Context, id string) (*domain.User, error) {
	data, err := c.client.Get(ctx, "user:"+id).Bytes()
	if err != nil {
		return nil, err
	}
	user := &domain.User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *UserCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, "user:"+id).Err()
}