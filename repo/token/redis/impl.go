package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepo struct {
	redisClient *redis.Client
}

func NewRefreshTokenRepository(rc *redis.Client) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		redisClient: rc,
	}
}

func (r *RefreshTokenRepo) Set(ctx context.Context, refreshToken string) error {
	return r.redisClient.Set(ctx, refreshToken, nil, 0).Err()
}

func (r *RefreshTokenRepo) Check(ctx context.Context, refreshToken string) error {
	_, err := r.redisClient.Get(ctx, refreshToken).Result()
	return err
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, refreshToken string) error {
	return r.redisClient.Del(ctx, refreshToken).Err()
}
