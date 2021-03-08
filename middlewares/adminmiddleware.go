package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware ...
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("ADMIN REQUEST")
		c.Next()

	}
}
