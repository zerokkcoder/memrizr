package service

import (
	"context"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"

	"github.com/google/uuid"
)

// 用户服务层结构体
type UserService struct {
	UserRepository model.UserRepository
}

// 用户服务层配置结构体
type USConfig struct {
	UserRepository model.UserRepository
}

// NewUserService 创建实例
func NewUserService(c *USConfig) model.UserService {
	return &UserService{
		UserRepository: c.UserRepository,
	}
}

// Get 实现 UserService 接口 Get 方法
func (s *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	return s.UserRepository.FindByID(ctx, uid)
}

// Signup 实现 UserService 接口 Signup 方法
func (s *UserService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	u.Password = pw
	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	return nil
}
