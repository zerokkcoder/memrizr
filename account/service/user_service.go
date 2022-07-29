package service

import (
	"context"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"

	"github.com/google/uuid"
)

// 用户服务层结构体
type userService struct {
	UserRepository model.UserRepository
}

// 用户服务层配置结构体
type USConfig struct {
	UserRepository model.UserRepository
}

// NewUserService 创建实例
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

// Get 实现 UserService 接口 Get 方法
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	return s.UserRepository.FindByID(ctx, uid)
}

// Signup 实现 UserService 接口 Signup 方法
func (s *userService) Signup(ctx context.Context, u *model.User) error {
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

// Signin 实现 UserService 接口 Signin 方法
func (s *userService) Signin(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}
	// 验证密码
	match, err := comparePasswords(uFetched.Password, u.Password)
	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*u = *uFetched
	return nil
}

// UpdateDetails 实现 UserService 接口 UpdateDetails 方法
func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	err := s.UserRepository.Update(ctx, u)
	if err != nil {
		return err
	}

	return nil
}
