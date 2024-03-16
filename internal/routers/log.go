package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FriedCoderZ/LogCollector-server/internal/database"
	"github.com/FriedCoderZ/LogCollector-server/internal/util"
	"github.com/gin-gonic/gin"
)

func SetupLog(router *gin.Engine) *gin.Engine {
	// 定义路由和处理函数
	router.POST("/logs/:uuid", InsertLogHandler)

	return router
}

func InsertLogHandler(c *gin.Context) {
	// 获取路由参数uuid的值
	uuid := c.Param("uuid")
	// 获取响应数据
	ciphertext, err := c.GetRawData()
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//数据库查询采集端对应信息
	collector, err := database.GetCollectorByUUID(uuid)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//数据库查询采集端对应AES Key
	aesKey, err := collector.GetAESKey()
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//AES解密
	plaintext, err := util.AesDecrypt(ciphertext, aesKey)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	var logsData []map[string]interface{}
	json.Unmarshal(plaintext, &logsData)
	fmt.Println(uuid)
	fmt.Println(collector)
	err = collector.AppendLogs(logsData)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
