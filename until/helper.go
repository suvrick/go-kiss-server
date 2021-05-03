package until

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
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

func GetMD5Hash(login, password string) string {
	hash := md5.Sum([]byte(login + password + time.Now().String()))
	return hex.EncodeToString(hash[:])
}
