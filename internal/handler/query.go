package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/FriedCoderZ/LogCollector-server/internal/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func QueryHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	query := c.Request.URL.Query()
	// 查询
	logs := []database.Log{}
	collection := database.GetDB().Collection("logs")
	bsons := append(buildLogQuery(query), bson.M{"uuid": uuid})
	querys := bson.M{"$and": bsons}
	cursor, err := collection.Find(context.TODO(), querys, options.Find())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = cursor.All(context.TODO(), &logs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(logs) == 0 {
		c.JSON(http.StatusOK, gin.H{"logs": []map[string]interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func buildLogQuery(query url.Values) []bson.M {
	logQuery := []bson.M{}
	for key, value := range query {
		// 解析条件和值
		value := value[0]
		intValue, _ := strconv.Atoi(value)
		cond, val := parseCondition(key)
		if val == "" {
			continue
		}
		// 构建查询子文档
		subQuery := bson.M{}
		switch cond {
		case "exact", "equal", "":
			subQuery["$or"] = []bson.M{
				{"data." + val: value},
				{"data." + val: intValue},
			}
		case "contains":
			subQuery["data."+val] = bson.M{"$regex": value}
		case "startswith":
			subQuery["data."+val] = bson.M{"$regex": "^" + value}
		case "endswith":
			subQuery["data."+val] = bson.M{"$regex": value + "$"}
		case "regular":
			subQuery["data."+val] = bson.M{"$regex": value + "$"}
		case "gt":
			subQuery["data."+val] = bson.M{"$gt": intValue}
		case "lt":
			subQuery["data."+val] = bson.M{"$lt": intValue}
		case "gte":
			subQuery["data."+val] = bson.M{"$gte": intValue}
		case "lte":
			subQuery["data."+val] = bson.M{"$lte": intValue}
		case "isnull":
			subQuery["data."+val] = nil
		case "in":
			values := strings.Split(value, ",")
			subQuery["data."+val] = bson.M{"$in": values}
		default:
			// 无效的条件，返回空查询
			return []bson.M{}
		}

		// 将子文档添加到主查询的切片中
		logQuery = append(logQuery, subQuery)
	}
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
