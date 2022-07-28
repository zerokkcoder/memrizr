package handler

import (
	"bytes"
	"encoding/json"
	"memrizr/model"
	"memrizr/model/apperrors"
	"memrizr/model/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup(t *testing.T) {
	// 设置 gin 模式
	gin.SetMode(gin.TestMode)

	// 测试没有邮箱和密码
	t.Run("Email and Password Required", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// ResponseRecorder 获取 http 响应
		rr := httptest.NewRecorder()

		// 路由
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// 创建请求体
		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		// 请求
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "signup")
	})

	// 无效邮箱测试用例
	t.Run("Invalid Email", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// ResponseRecorder 获取 http 响应
		rr := httptest.NewRecorder()

		// 路由
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// 创建请求体
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bo",
			"password": "avalidpassword123",
		})
		assert.NoError(t, err)

		// 请求
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "signup")
	})

	// 密码太短测试用例
	t.Run("Password too short", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// ResponseRecorder 获取 http 响应
		rr := httptest.NewRecorder()

		// 路由
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// 创建请求体
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"password": "inval",
		})
		assert.NoError(t, err)

		// 请求
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "signup")
	})

	// 密码太长测试用例
	t.Run("Password too long", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// ResponseRecorder 获取 http 响应
		rr := httptest.NewRecorder()

		// 路由
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// 创建请求体
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"password": "invalkadsfjasdfkj;askldfj;askldfj;asdfiuerueuuuuuudfjasdfasdkfjj",
		})
		assert.NoError(t, err)

		// 请求
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "signup")
	})

	// 通过用户服务注册方法返沪错误 测试用例
	t.Run("Error calling UserService", func(t *testing.T) {
		u := &model.User{
			Email:    "1111@qq.com",
			Password: "123456789",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), u).Return(apperrors.NewConflict("User Already Exists", u.Email))

		// ResponseRecorder 获取 http 响应
		rr := httptest.NewRecorder()

		// 路由
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// 创建请求体
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// 请求
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 500, rr.Code)
		mockUserService.AssertNotCalled(t, "signup")
	})
}
