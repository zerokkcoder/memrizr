package model

// TokenPair 返回 idtoken 和 refreshToken
type TokenPair struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}


