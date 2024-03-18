package routers

import (
	"github.com/FriedCoderZ/LogCollector-server/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupQuery(router *gin.Engine) *gin.Engine {
	router.GET("/query/:uuid", handler.QueryHandler)
	return router
}
