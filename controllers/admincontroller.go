package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
func NewAdminController(r *gin.Engine, us *services.UserService) {
	ctrl := &AdminController{
		router:      r,
		userService: us,
	}

	admin := ctrl.router.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	admin.GET("/get/:id", ctrl.getUserHandler)
	admin.GET("/setlimit/:id/:limit", ctrl.setLimitBotToUserHandler)
	admin.GET("/setdate/:id/:mounth", ctrl.setDateBotToUserHandler)
}

func (ctrl *AdminController) setLimitBotToUserHandler(c *gin.Context) {
	id := c.Param("id")
	limitStr := c.Param("limit")

	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	user, err := ctrl.userService.FindByID(userID)

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
	mounthStr := c.Param("mounth")

	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	mounth, err := strconv.Atoi(mounthStr)
	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	user, err := ctrl.userService.FindByID(userID)

	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	newTime := time.Now().AddDate(0, mounth, 0)
	user.Date = newTime.Format("2006-01-02")
	ctrl.userService.UpdateUser(user)

	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}

func (ctrl *AdminController) getUserHandler(c *gin.Context) {
	id := c.Param("id")

	userID, _ := strconv.Atoi(id)

	user, err := ctrl.userService.FindByID(userID)
	until.WriteResponse(c, 200, gin.H{
		"user": user,
	}, err)
}
