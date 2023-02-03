package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var allowOrigins = []string{
	"https://api.wxwind.top",
	"https://www.wxwind.top",
}

func IsContain[T string | int](items []T, item T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Cors() gin.HandlerFunc {

	return func(c *gin.Context) {
		method := c.Request.Method

		//是cors请求
		origin := c.Request.Header.Get("Origin")
		if origin != "" && IsContain(allowOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
			//允许浏览器发送的头
			c.Header("Access-Control-Allow-Headers", "Content-Type")
			//允许浏览器拿到的头
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			//是否允许cookies, authorization headers 或 TLS client certificates
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//是非简单请求的预检请求，直接返回204，不做后续处理
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
