package router

import (
	"common/config"
	"common/rpc"
	"gateway/api"
	"gateway/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 注册路由
func RegisterRouter() *gin.Engine {
	if config.Conf.Log.Level == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化grpc的client gate是做为grpc的客户端 去调用user grpc服务
	rpc.Init()
	engine := gin.Default()
	engine.Use(auth.Cors())
	userHandler := api.NewUserHandler()
	engine.POST("/register", userHandler.Register)

	return engine
}
