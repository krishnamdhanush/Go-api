package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krishnamdhanush/web-dawn/configs"
)

func AlbumCacheById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		id, check := c.Params.Get("id")

		if check {
			value, err := configs.RedisClient.Get(ctx, fmt.Sprintf("albums:%s", id)).Result()
			if err == nil {
				c.IndentedJSON(http.StatusFound, value)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func AlbumCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		value, err := configs.RedisClient.Get(ctx, "albums:all").Result()
		if err == nil {
			c.AbortWithStatusJSON(http.StatusFound, value)
			return
		}
		c.Next()
	}
}
