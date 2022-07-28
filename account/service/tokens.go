package service

import (
	"crypto/rsa"
	"log"
	"memrizr/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// IDTokenCustomClaims 自定义jwt claims
type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

// 生成 token
func generateIDToken(u *model.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + 60*15 // 15分钟过期实践

	claims := IDTokenCustomClaims{
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

// RefreshToken 刷新Token 结构体
type RefreshToken struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

// RefreshTokenCustomClaims 刷新token的自定义jwt claims
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// 生成刷新token
func generateRefreshToken(uid uuid.UUID, key string) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.AddDate(0, 0, 3) // 3天过期时间
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Faild to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
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

	return &RefreshToken{
		SS:        ss,
		ID:        tokenID.String(),
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}
