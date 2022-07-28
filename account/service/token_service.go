package service

import (
	"context"
	"crypto/rsa"
	"log"
	"memrizr/model"
	"memrizr/model/apperrors"
)

// TokenService Token服务层
type tokenService struct {
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TSConfig Token服务层配置结构体
type TSConfig struct {
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenService 实例化TokenService
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		PrivateKey:            c.PrivateKey,
		PublicKey:             c.PublicKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewTokenPairFromUser 实现方法
func (s *tokenService) NewTokenPairFromUser(ctx context.Context, u *model.User, prevIDToken string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivateKey, s.IDExpirationSecs)
	if err != nil {
		log.Printf("Error generateing idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)
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
