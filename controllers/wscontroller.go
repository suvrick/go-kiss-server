package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/services"
)

type WsController struct {
	router      *gin.Engine
	userService *services.UserService
}

func NewWsController(r *gin.Engine, u_service *services.UserService) *WsController {
	ctrl := &WsController{
		router:      r,
		userService: u_service,
	}

	ctrl.router.Use(middlewares.AuthMiddleware())
	ctrl.router.GET("/ws", ctrl.acceptHandler)

	return ctrl
}

func (ctrl *WsController) acceptHandler(c *gin.Context) {

}
