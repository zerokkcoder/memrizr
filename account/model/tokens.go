package model

import "github.com/google/uuid"

// RefreshToken 存储 token 属性
type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refreshToken"`
}

// IDToken 存储 token 属性
type IDToken struct {
	SS string `json:"idToken"`
}

// TokenPair 返回 idtoken 和 refreshToken
type TokenPair struct {
	IDToken
	RefreshToken
}
