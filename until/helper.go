package until

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/model"
)

// Response ...
type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

// WriteResponse ...
func WriteResponse(c *gin.Context, code int, data gin.H, err error) {
	var msgErr string
	if err != nil {
		msgErr = err.Error()
	}

	response := &Response{
		Code:  code,
		Data:  data,
		Error: msgErr,
	}

	c.JSON(code, response)
	c.Abort()
}

// GetUserFromContext ...
func GetUserFromContext(c *gin.Context) (string, model.User, error) {

	u, ok := c.Get("user")
	if !ok {
		return "", model.User{}, errors.ErrNotAuthenticated
	}

	user := u.(model.User)
	return strconv.Itoa(user.ID), user, nil
}
