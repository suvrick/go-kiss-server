package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/until"
)

// AdminMiddleware ...
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		u, ok := c.Get("user")
		if !ok {
			until.WriteResponse(c, 403, nil, errors.New("not forbidden"))
			c.Abort()
			return
		}

		user := u.(model.User)
		if user.Email != "s@mail.com" {
			until.WriteResponse(c, 403, nil, errors.New("not forbidden"))
			c.Abort()
			return
		}

		c.Next()

	}
}
