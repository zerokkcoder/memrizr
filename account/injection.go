package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"memrizr/handler"
	"memrizr/repository"
	"memrizr/service"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 初始化 处理器
// 注入存储层
// 注入服务层
// 注入处理层
func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	// 存储层
	userRepository := repository.NewUserRepository(d.DB)
	// 服务层
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})
	// 加载 rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %W\n", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %W\n", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %W\n", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %W\n", err)
	}

	// 从 env 中加载 refresh token secret
	refreshSecret := os.Getenv("REFRESH_SECRET")

	tokenService := service.NewTokenService(&service.TSConfig{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		RefreshSecret: refreshSecret,
	})

	// 路由器
	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R:            router,
		UserService:  userService,
		TokenService: tokenService,
	})

	return router, nil
}
