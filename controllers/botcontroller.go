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
func NewBotController(r *gin.Engine, bs *services.BotService, us *services.UserService) {

	ctrl := &BotController{
		router:      r,
		botService:  bs,
		userService: us,
	}

	group := ctrl.router.Group("bots")

	group.Use(middlewares.AuthMiddleware())

	group.GET("/all", ctrl.allHandler)
	group.POST("/add", ctrl.addHandler)
	group.GET("/remove/:botID", ctrl.removeHandler)

}

func (ctrl *BotController) addHandler(c *gin.Context) {

	type U struct {
		URL string `json:"url"`
	}

	postForm := &U{}
	if err := c.ShouldBindJSON(postForm); err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	bot, err := ctrl.botService.Add(c, postForm.URL)

	until.WriteResponse(c, 200, gin.H{
		"bot": bot,
	}, err)
}

func (ctrl *BotController) allHandler(c *gin.Context) {

	bots, err := ctrl.botService.All(c)

	if err != nil {
		until.WriteResponse(c, 200, nil, err)
		return
	}

	until.WriteResponse(c, 200, gin.H{
		"count": len(bots),
		"bots":  bots,
	}, nil)
}

func (ctrl *BotController) removeHandler(c *gin.Context) {

	err := ctrl.botService.Delete(c)
	ok := true
	if err != nil {
		ok = false
	}

	until.WriteResponse(c, 200, gin.H{
		"result": ok,
	}, err)
}
