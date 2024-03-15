package main

import (
	"github.com/FriedCoderZ/LogCollector-server/internal/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin路由器
	router := gin.Default()

	// 定义路由和处理函数
	routers.SetupCollector(router)
	// routers.SetupQuery(router)
	// 启动HTTP服务器，默认监听在localhost:8080
	router.Run()
}
