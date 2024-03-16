package routers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/FriedCoderZ/LogCollector-server/internal/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SetupQuery(router *gin.Engine) *gin.Engine {
	// 定义路由和处理函数
	router.GET("/query/:uuid", handleQuery)

	return router
}

func handleQuery(c *gin.Context) {
	uuid := c.Param("uuid")
	query := c.Request.URL.Query()
	// 查询
	collector := []database.Collector{}
	collection := database.GetDB().Collection("collectors")
	bsons := append(buildLogQuery(query), bson.M{"uuid": uuid})
	fmt.Println(bsons)
	querys := bson.M{"$and": bsons}
	cursor, err := collection.Find(context.TODO(), querys)
	// cursor, err := collection.Find(context.TODO(), bson.M{"uuid": uuid, "logs": buildLogQuery(query)})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = cursor.All(context.TODO(), &collector)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(collector) == 0 {
		c.JSON(http.StatusOK, gin.H{"logs": []map[string]interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": collector[0].Logs})
}

func buildLogQuery(query url.Values) []bson.M {
	logQuery := []bson.M{}
	for key, value := range query {
		// 解析条件和值
		value := value[0]
		cond, val := parseCondition(key)
		fmt.Println(cond)
		if val == "" {
			continue
		}
		// 构建查询子文档
		subQuery := bson.M{}
		switch cond {
		case "exact", "iexact", "":
			subQuery["logs."+val] = value
		case "contains", "icontains":
			subQuery["logs."+val] = bson.M{"$regex": value}
		case "startswith", "istartswith":
			subQuery["logs."+val] = bson.M{"$regex": "^" + value}
		case "endswith", "iendswith":
			subQuery["logs."+val] = bson.M{"$regex": value + "$"}
		case "range":
			parts := strings.Split(value, "_")
			if len(parts) != 2 {
				continue
			}
			start := parts[0]
			end := parts[1]
			subQuery["logs."+val] = bson.M{"$gte": start, "$lte": end}
		case "gt":
			subQuery["logs."+val] = bson.M{"$gt": value}
		case "lt":
			subQuery["logs."+val] = bson.M{"$lt": value}
		case "gte":
			subQuery["logs."+val] = bson.M{"$gte": value}
		case "lte":
			subQuery["logs."+val] = bson.M{"$lte": value}
		case "isnull":
			subQuery["logs."+val] = nil
		case "in":
			values := strings.Split(value, ",")
			subQuery["logs."+val] = bson.M{"$in": values}
		default:
			// 无效的条件，返回空查询
			return []bson.M{}
		}

		// 将子文档添加到主查询的切片中
		fmt.Println(logQuery)
		logQuery = append(logQuery, subQuery)
	}
	fmt.Println(logQuery)
	return logQuery
}

func parseCondition(condition string) (string, string) {
	// 分离条件类型和字段名
	parts := strings.Split(condition, "_")
	if len(parts) == 0 {
		return "", ""
	} else if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[1], parts[0]
}
