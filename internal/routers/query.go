package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupQuery(router *gin.Engine) *gin.Engine {
	// 定义路由和处理函数
	router.GET("/query/:id", func(c *gin.Context) {
		// 获取路由参数id的值
		id := c.Param("id")

		// 获取所有的查询参数
		queryParams := c.Request.URL.Query()

		// 存储参数的map
		params := make(map[string]string)

		// 遍历查询参数，将参数存储到map中
		for key, value := range queryParams {
			params[key] = value[0]
		}

		// 返回存储的参数map和路由参数id的值
		c.JSON(http.StatusOK, gin.H{"id": id, "params": params})
	})

	return router
}
