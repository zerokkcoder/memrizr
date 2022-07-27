package model

import (
	"context"

	"github.com/google/uuid"
)

// UserService 用户处理服务
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
}

// UserRepository 用户存储服务
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
}
