package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/until"
	"gorm.io/gorm"
)

// AdminController ...
type AdminController struct {
	router               *gin.Engine
	userService          *services.UserService
	kissUserService      *services.AutoKissService
	stateDownloadService *services.StateDownloadService
}

// NewAdminController ...
func NewAdminController(r *gin.Engine, us *services.UserService, kus *services.AutoKissService, sds *services.StateDownloadService) *AdminController {
	ctrl := &AdminController{
		router:               r,
		userService:          us,
		kissUserService:      kus,
		stateDownloadService: sds,
	}

	admin := ctrl.router.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	admin.GET("/get/:id", ctrl.getUserHandler)
	admin.GET("/setlimit/:id/:limit", ctrl.setLimitBotToUserHandler)
	admin.GET("/setdate/:id/:mounth", ctrl.setDateBotToUserHandler)

	kiss := admin.Group("/kiss")
	kiss.GET("/all", ctrl.allKissUserHandler)
	kiss.GET("/auth/:userID", ctrl.authKissUserHandler)
	kiss.GET("/get/:userID", ctrl.getKissUserHandler)
	kiss.GET("/down", ctrl.allDownloadHandler)

	return ctrl
}

func (ctrl *AdminController) allDownloadHandler(c *gin.Context) {
	states, _ := ctrl.stateDownloadService.AllDownloadState()
	c.JSON(200, gin.H{
		"count":  len(states),
		"states": states,
	})
}

func (ctrl *AdminController) authKissUserHandler(c *gin.Context) {

	id := c.Param("userID")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, gin.H{"result": "fail", "error": err.Error()})
	}

	user, err := ctrl.kissUserService.FindByIDKissUser(userID)

	if err == gorm.ErrRecordNotFound {
		err = ctrl.kissUserService.AddKissUser(userID)
	}

	if user.UserID == 0 {
		user = &model.KissUser{
			UserID: userID,
		}
	}

	user.IsTrial = false

	err = ctrl.kissUserService.UpdateKissUser(user)
	if err != nil {
		c.JSON(200, gin.H{"result": "fail", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": "ok", "user": user})

}

func (ctrl *AdminController) getKissUserHandler(c *gin.Context) {

	id := c.Param("userID")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, gin.H{"result": "fail"})
	}

	user, err := ctrl.kissUserService.FindByIDKissUser(userID)

	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"user":   nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"result": "ok",
		"user":   user,
	})

}

func (ctrl *AdminController) allKissUserHandler(c *gin.Context) {

	users, err := ctrl.kissUserService.AllKissUser()
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"users":  nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"result": "ok",
		"count":  len(users),
		"users":  users,
	})
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
	monthStr := c.Param("mounth")

	userID, err := strconv.Atoi(id)
	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	mounth, err := strconv.Atoi(monthStr)
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
