package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/until"
)

// BotController ...
type BotController struct {
	router      *gin.Engine
	botService  *services.BotService
	userService *services.UserService
}

// NewBotController ...
func NewBotController(r *gin.Engine, b_service *services.BotService, u_service *services.UserService) {

	ctrl := &BotController{
		router:      r,
		botService:  b_service,
		userService: u_service,
	}

	bots := ctrl.router.Group("bots", middlewares.AuthMiddleware())
	{
		bots.GET("/all", ctrl.allHandler)
		bots.POST("/add", ctrl.addHandler)
		bots.GET("/update/:botID", ctrl.updateHandler)
		bots.GET("/remove/:botID", ctrl.removeHandler)
	}

}

func (ctrl *BotController) addHandler(c *gin.Context) {

	type FormData struct {
		URL string `json:"url"`
	}

	data := &FormData{}
	if err := c.ShouldBindJSON(data); err != nil {
		until.WriteResponse(c, 201, nil, errors.ErrInvalidParam)
		return
	}

	// urls := strings.Split(data.URL, "\n")

	bot, err := ctrl.botService.Add(c, data.URL)

	until.WriteResponse(c, 200, gin.H{
		"bot": bot,
	}, err)

}

func (ctrl *BotController) updateHandler(c *gin.Context) {

	botUID := c.Param("botID")
	bot, err := ctrl.botService.UpdateByID(c, botUID)
	until.WriteResponse(c, 200, gin.H{
		"bot": bot,
	}, err)

}

func (ctrl *BotController) allHandler(c *gin.Context) {

	bots, err := ctrl.botService.All(c)

	if err != nil {
		until.WriteResponse(c, 201, gin.H{
			"count": len(bots),
			"bots":  bots,
		}, err)
	}

	until.WriteResponse(c, 200, gin.H{
		"count": len(bots),
		"bots":  bots,
	}, nil)
}

func (ctrl *BotController) removeHandler(c *gin.Context) {

	if err := ctrl.botService.Delete(c); err != nil {
		until.WriteResponse(c, 201, gin.H{
			"result": "fail",
		}, err)
		return
	}

	until.WriteResponse(c, 200, gin.H{
		"result": "ok",
	}, nil)

}
