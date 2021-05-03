package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/until"
)

// AuthMiddleware ...
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := session.GetUser(c)

		if user == nil || user.ID == 0 {
			until.WriteResponse(c, 401, gin.H{
				"result": "fail",
			}, errors.ErrNotAuthenticated)
			c.Abort()
			return
		}

		c.Next()
	}
}
