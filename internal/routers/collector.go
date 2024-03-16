package routers

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/FriedCoderZ/LogCollector-server/internal/database"
	"github.com/FriedCoderZ/LogCollector-server/internal/util"
	"github.com/gin-gonic/gin"
)

func SetupCollector(router *gin.Engine) *gin.Engine {
	// 定义路由和处理函数
	router.POST("/collector", collectorHandler)

	return router
}

func collectorHandler(c *gin.Context) {
	requestBody, err := c.GetRawData()
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	encryptedKey, _ := base64.StdEncoding.DecodeString(string(requestBody))
	privateKey, err := os.ReadFile("/ouryun/LogCollector-server/privateKey.pem")
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// 解密获得AES秘钥
	aesKey, err := util.RSADecrypt(encryptedKey, privateKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// 将AES密钥存储到数据库以生成一条数据，获取ID
	collector, err := database.CreateCollector(aesKey)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// 返回ID作为响应
	c.JSON(http.StatusOK, gin.H{"uuid": collector.UUID})
}
