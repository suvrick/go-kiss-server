package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/until"
)

// UserController ...
type UserController struct {
	router      *gin.Engine
	userService *services.UserService
}

// NewUserController ...
func NewUserController(r *gin.Engine, s *services.UserService) {

	ctrl := &UserController{
		router:      r,
		userService: s,
	}

	user := ctrl.router.Group("/user")

	user.POST("/login", ctrl.loginHandler)
	user.POST("/register", ctrl.registerHandler)

	self := ctrl.router.Group("/self")
	self.Use(middlewares.AuthMiddleware())
	self.GET("/get", ctrl.getUserHandler)

}

func (ctrl *UserController) registerHandler(c *gin.Context) {
	type L struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	login := L{}

	if err := c.ShouldBindJSON(&login); err != nil {
		until.WriteResponse(c, 200, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	if len(login.Login) == 0 || len(login.Password) == 0 {
		until.WriteResponse(c, 200, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	id, err := ctrl.userService.Create(login.Login, login.Password)
	if err != nil {
		until.WriteResponse(c, 200, gin.H{
			"result": "fail",
		}, errors.ErrIncorrectEmailOrPassword)
		return
	}

	until.WriteResponse(c, 200, gin.H{
		"result": "ok",
		"id":     id,
	}, nil)
	return
}

func (ctrl *UserController) loginHandler(c *gin.Context) {

	type L struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	login := L{}

	if err := c.ShouldBindJSON(&login); err != nil {
		until.WriteResponse(c, 200, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	if len(login.Login) == 0 || len(login.Password) == 0 {
		until.WriteResponse(c, 200, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	user, err := ctrl.userService.Login(c, login.Login, login.Password)

	if err != nil {
		until.WriteResponse(c, 403, gin.H{
			"result": "ok",
			"user":   user,
		}, errors.ErrIncorrectEmailOrPassword)
		return
	}

	until.WriteResponse(c, 200, gin.H{
		"result": "ok",
		"user":   user,
	}, nil)
}

func (ctrl *UserController) getUserHandler(c *gin.Context) {

	_, user, err := until.GetUserFromContext(c)
	fmt.Println(user)
	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}
