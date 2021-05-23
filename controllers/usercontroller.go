package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/until"
)

// UserController ...
type UserController struct {
	router      *gin.Engine
	userService *services.UserService
}

/*



	HTTP Status code:
	200 - OK
	400 - BadRequest
	401 - Unauthorized
	403 - Not Forbidden




*/

// NewUserController ...
func NewUserController(r *gin.Engine, uService *services.UserService) {

	ctrl := &UserController{
		router:      r,
		userService: uService,
	}

	ctrl.router.GET("/logout", ctrl.logoutHandler)

	user := ctrl.router.Group("/user")
	{
		user.GET("/logout", ctrl.logoutHandler)
		user.POST("/login", ctrl.loginHandler)
		user.POST("/register", ctrl.registerHandler)
	}

	self := ctrl.router.Group("/user", middlewares.AuthMiddleware())
	{
		self.GET("/get", ctrl.getUserHandler)
	}

}

func (ctrl *UserController) registerHandler(c *gin.Context) {

	type FormData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data := FormData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		until.WriteResponse(c, 400, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	if len(data.Email) == 0 || len(data.Password) == 0 {
		until.WriteResponse(c, 400, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	id, err := ctrl.userService.Register(data.Email, data.Password)

	if err != nil {
		until.WriteResponse(c, 401, gin.H{
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

	type FormData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data := FormData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		until.WriteResponse(c, 400, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	if len(data.Email) == 0 || len(data.Password) == 0 {
		until.WriteResponse(c, 400, gin.H{
			"result": "fail",
		}, errors.ErrInvalidParam)
		return
	}

	user, err := ctrl.userService.Login(data.Email, data.Password)

	if err != nil {
		until.WriteResponse(c, 401, gin.H{
			"result": "fail",
			"user":   user,
		}, errors.ErrIncorrectEmailOrPassword)
		return
	}

	ctrl.setCookie(c, user)

	until.WriteResponse(c, 200, gin.H{
		"result": "ok",
		"user":   user,
	}, nil)
}

func (ctrl *UserController) logoutHandler(c *gin.Context) {

	user := session.GetUser(c)
	if user == nil {
		until.WriteResponse(c, 401, nil, nil)
		return
	}

	ctrl.deleteCookie(c, user)
	if err := ctrl.userService.UpdateUser(user); err != nil {
		//...
	}

	until.WriteResponse(c, 401, nil, nil)
}

// getUserHandler возращает текущего пользователя на ружу
// TODO: перенести в wss ?
func (ctrl *UserController) getUserHandler(c *gin.Context) {

	user := session.GetUser(c)

	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, nil)
}

func (ctrl *UserController) setCookie(c *gin.Context, user *model.User) {

	//Ну тут приметивный хеш
	user.Token = until.GetMD5Hash(user.Email, user.Password)

	host := strings.Split(c.Request.Host, ":")[0] //отрезаем порт (нужен для локалки)
	c.SetCookie("token", user.Token, 60*60*30*24, "/", host, false, false)

	// TODO:Обработать ошибку
	err := ctrl.userService.UpdateUser(user)
	if err != nil {
		return
	}
}

func (ctrl *UserController) deleteCookie(c *gin.Context, user *model.User) {
	delete(session.Accounts, user.Token)
	user.Token = ""
	host := strings.Split(c.Request.Host, ":")[0]
	c.SetCookie("token", user.Token, -1, "/", host, false, false)
}
