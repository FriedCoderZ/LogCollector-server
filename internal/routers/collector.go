package routers

import (
	"github.com/FriedCoderZ/LogCollector-server/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupCollector(router *gin.Engine) *gin.Engine {
	router.POST("/collector", handler.CollectorHandler)
	return router
}
