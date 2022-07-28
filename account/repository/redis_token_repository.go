package repository

import (
	"context"
	"fmt"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken 存储token
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

// DeleteRefreshToken 删除token
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, prevTokenID)
	if err := r.Redis.Del(ctx, key); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, prevTokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}
