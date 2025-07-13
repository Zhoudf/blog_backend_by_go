package middleware

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 错误恢复中间件，捕获panic并返回友好错误
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈信息
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				log.Printf("panic recovered:\n%s\n%s", err, buf[:n])

				// 返回500错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})

				// 终止请求处理链
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}