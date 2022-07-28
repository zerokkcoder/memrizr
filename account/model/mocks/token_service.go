package mocks

import (
	"context"
	"memrizr/model"

	"github.com/stretchr/testify/mock"
)

// MockTokenSerive 模拟Token服务
type MockTokenSerive struct {
	mock.Mock
}

// NewTokenPairFromUser 模拟生成token
func (m *MockTokenSerive) NewTokenPairFromUser(ctx context.Context, u *model.User, prevIDToken string) (*model.TokenPair, error) {
	ret := m.Called(ctx, u, prevIDToken)

	var r0 *model.TokenPair
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.TokenPair)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
