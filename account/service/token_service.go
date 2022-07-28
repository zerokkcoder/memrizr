package service

import (
	"context"
	"crypto/rsa"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"
)

// TokenService Token服务层
type TokenService struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// TSConfig Token服务层配置结构体
type TSConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// NewTokenService 实例化TokenService
func NewTokenService(c *TSConfig) model.TokenService {
	return &TokenService{
		PrivateKey:    c.PrivateKey,
		PublicKey:     c.PublicKey,
		RefreshSecret: c.RefreshSecret,
	}
}

// NewTokenPairFromUser 实现方法
func (s *TokenService) NewTokenPairFromUser(ctx context.Context, u *model.User, prevIDToken string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivateKey)
	if err != nil {
		log.Printf("Error generateing idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret)
	if err != nil {
		log.Printf("Error generateing refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// TODO 保存 refresh token

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
