package routers

import (
	"github.com/FriedCoderZ/LogCollector-server/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupLog(router *gin.Engine) *gin.Engine {
	router.POST("/logs/:uuid", handler.InsertLogHandler)
	return router
}
