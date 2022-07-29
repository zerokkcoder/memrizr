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

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, prevTokenID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/TokenID: %s/%s does not exists\n", userID, prevTokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

// DeleteUserRefreshTokens 退出删除用户token
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s", iter.Val())
			failCount++
		}
	}

	// 检验最后的值
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
