# Golang+Ract全栈APP项目实战


## 使用到的库
- `go get -u github.com/gin-gonic/gin` http服务器及理路由
- `go get github.com/cespare/reflex` 是应用程序进行实时重新加载
- `go get github.com/stretchr/testify` 测试包
- `go get github.com/google/uuid` 创建 uuid 包



## 命令列表
启动 `docker-compose`:
```
$ make up
```
停止 `docker-compose`:
```
$ make down
```
生成公私密钥文件：
```
$ create-keypair
```
创建迁移文件(需安装 `migrate`)：
```
$ make migrate-down
```
迁移表:
```
$ make migtate-up [N=number]
```
删除表:
```
$ make migtate-up [N=number]
```
强制迁移:
```
$ make migrate-force VERSION=number
```

## 图文教程列表
- [01-Docker重载安装Go服务器](https://dev.to/jacobsngoodwin/full-stack-memory-app-01-setup-go-server-with-reload-in-docker-62n)
- [02-使用 gin 创建路由处理器](https://dev.to/jacobsngoodwin/02-creating-route-handlers-in-gin-4f3j)
- [03-应用架构](https://dev.to/jacobsngoodwin/03-application-architecture-5jk)
- [04-使用 Testify Mock 测试 Gin HTTP 处理器](https://dev.to/jacobsngoodwin/04-testing-first-gin-http-handler-9m0)
- [05-测试GO Account应用服务层方法](https://dev.to/jacobsngoodwin/05-testing-a-service-layer-method-in-go-account-application-13a6)
- [06-使用 gin 创建注册处理器 - 绑定数据](https://dev.to/jacobsngoodwin/creating-signup-handler-in-gin-binding-data-3kb5)
- [07-使用 gin 创建注册处理器 - 生成Token](https://dev.to/jacobsngoodwin/07-completing-signup-handler-in-gin-token-creation-1ikc)
- [08-实现注册的服务层和存储层](https://dev.to/jacobsngoodwin/08-implement-signup-in-service-and-repository-layers-4coe)
- [09-Token 生成](https://dev.to/jacobsngoodwin/09-token-creation-gjh)
- [10-依赖注入及应用演示](https://dev.to/jacobsngoodwin/10-dependency-injection-and-app-demo-1pj5)
- [11-Account API清理和修复](https://dev.to/jacobsngoodwin/11-cleanup-fixes-2b18)
- [12-Redis 保存 RefreshToken](https://dev.to/jacobsngoodwin/12-store-refresh-tokens-in-redis-1k5d)
- [13-Gin 超时处理中间件](https://dev.to/jacobsngoodwin/13-gin-handler-timeout-middleware-4bhg)
- [14-添加 Signin 登录路由处理器](https://dev.to/jacobsngoodwin/14-add-signin-handler-32be)
- [15-服务层、存储层实现 Signin](https://dev.to/jacobsngoodwin/15-add-signin-to-service-and-repository-layers-5mg)