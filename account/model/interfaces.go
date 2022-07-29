package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService 用户处理服务
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
	Signin(ctx context.Context, u *User) error
}

// UserRepository 用户存储服务
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// TokenService Token服务接口
type TokenService interface {
	NewTokenPairFromUser(ctx context.Context, u *User, prevIDToken string) (*TokenPair, error)
}

// TokenRepository Token存储接口
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
}
