package middlewares

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		w := c.Writer
		//r := c.Request

		// js-ajax handling cross-domain issues
		//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin, X-Requested-With, Content-Type, Accept, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
	// return cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:4200"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "GET"},
	// 	AllowHeaders:     []string{"Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	// 	AllowCredentials: true,
	// })

	// return func(c *gin.Context) {

	// 	c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
	// 	c.Header("Access-Control-Allow-Methods", "*")
	// 	c.Header("Access-Control-Allow-Headers", "*")
	// 	c.Header("Content-Type", "application/json")
	// 	// c.Writer.Header().Add("Access-Control-Allow-Origin", "origin")
	// 	// //c.Writer.Header().Add("Access-Control-Allow-Origin", "https://suvricksoft.ru")
	// 	// //c.Writer.Header().Add("Access-Control-Allow-Origin", "http://localhost:4200")
	// 	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
	// 	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}

	// 	c.Next()
	// }
}
