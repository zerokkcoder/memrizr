package service

import (
	"crypto/rsa"
	"fmt"
	"log"
	"memrizr/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// idTokenCustomClaims 自定义jwt claims
type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

// 生成 token
func generateIDToken(u *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := idTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Println("Faild to sign id token string")
		return "", err
	}

	return ss, nil
}

// refreshTokenData 刷新Token 结构体
type refreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

// refreshTokenCustomClaims 刷新token的自定义jwt claims
type refreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// 生成刷新token
func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Faild to generate refresh token ID")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Faild to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SS:        ss,
		ID:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

// validateIDToken 验证 token
func validateIDToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// 现在我们只返回错误并处理服务级别的登录
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

// validateRefreshToken 验证 refreshToken
func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	// 现在我们只返回错误并处理服务级别的登录
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Refresh token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}

	return claims, nil
}
