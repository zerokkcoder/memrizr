package model

import (
	"context"

	"github.com/google/uuid"
)

// UserService 用户处理服务
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
}

// UserRepository 用户存储服务
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
}

// TokenService Token服务接口
type TokenService interface {
	NewTokenPairFromUser(ctx context.Context, u *User, prevIDToken string) (*TokenPair, error)
}
