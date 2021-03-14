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

		token, err := c.Cookie("token")

		if err != nil {

			until.WriteResponse(c, 403, gin.H{
				"result": "fail",
			}, errors.ErrNotAuthenticated)
			c.Abort()
			return
		}

		c.Set("token", token)

		user, ok := session.Accounts[token]

		if !ok {
			until.WriteResponse(c, 403, gin.H{
				"result": "fail",
			}, errors.ErrNotAuthenticated)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()

	}
}
