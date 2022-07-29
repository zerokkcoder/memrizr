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
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TSConfig Token服务层配置结构体
type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenService 实例化TokenService
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
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

	// 保存 refresh token
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	if prevIDToken != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevIDToken); err != nil {
			log.Printf("could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevIDToken)
		}
	}

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}

// ValidateIDToken 验证 token
func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, s.PublicKey)
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}
	return claims.User, nil
}
