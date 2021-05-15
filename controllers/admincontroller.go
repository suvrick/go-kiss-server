package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/until"
)

// AdminController ...
type AdminController struct {
	router      *gin.Engine
	userService *services.UserService
}

// NewAdminController ...
func NewAdminController(r *gin.Engine, u_service *services.UserService) *AdminController {
	ctrl := &AdminController{
		router:      r,
		userService: u_service,
	}

	admin := ctrl.router.Group("/admin")

	admin.Use(middlewares.AuthMiddleware()).Use(middlewares.AdminMiddleware())

	botovod := admin.Group("/botovod")
	{
		botovod.GET("/all", ctrl.getAllUser)
		botovod.GET("/get/:id", ctrl.getUserHandler)
		botovod.GET("/setlimit/:id/:limit", ctrl.setLimitBotToUserHandler)
		botovod.GET("/setdate/:id/:mounth", ctrl.setDateBotToUserHandler)
	}

	return ctrl
}

func (ctrl *AdminController) setLimitBotToUserHandler(c *gin.Context) {
	id := c.Param("id")
	limitStr := c.Param("limit")

	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	user, err := ctrl.userService.FindUserByID(userID)

	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	user.Limit = limit
	ctrl.userService.UpdateUser(user)

	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}

func (ctrl *AdminController) setDateBotToUserHandler(c *gin.Context) {

	id := c.Param("id")
	monthStr := c.Param("mounth")

	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	mounth, err := strconv.Atoi(monthStr)
	if err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	user, err := ctrl.userService.FindUserByID(userID)

	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	newTime := time.Now().AddDate(0, mounth, 0)
	user.Date = newTime.Format(until.TIME_FORMAT)

	ctrl.userService.UpdateUser(user)

	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}

func (ctrl *AdminController) getUserHandler(c *gin.Context) {

	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	user, err := ctrl.userService.FindUserByID(userID)

	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}

func (ctrl *AdminController) getAllUser(c *gin.Context) {
	users, err := ctrl.userService.AllUser()
	until.WriteResponse(c, 200, gin.H{
		"count": len(users),
		"users": users,
	}, err)
}
