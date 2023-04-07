package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		c.Set("page", page)
		c.Set("limit", limit)
		c.Next()
	}
}
