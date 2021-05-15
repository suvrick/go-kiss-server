package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/until"
)

// ProxyController ...
type ProxyController struct {
	router       *gin.Engine
	proxyService *services.ProxyService
}

// NewProxyController ...
func NewProxyController(r *gin.Engine, s *services.ProxyService) *ProxyController {
	ctrl := &ProxyController{
		router:       r,
		proxyService: s,
	}

	proxy := ctrl.router.Group("/proxy", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		proxy.POST("/add", ctrl.addHandler)
		proxy.GET("/all", ctrl.allHandler)
		proxy.GET("/free", ctrl.freeHandler)
		proxy.GET("/clear", ctrl.clearAllHandler)
	}

	return ctrl
}

// allHandler
// Получем все прокси из базы
func (ctrl *ProxyController) allHandler(c *gin.Context) {
	prxs, err := ctrl.proxyService.All()

	until.WriteResponse(c, 200, gin.H{
		"count":   len(prxs),
		"proxies": prxs,
	}, err)
}

// addHandler
// Добавляем прокси в бд
func (ctrl *ProxyController) addHandler(c *gin.Context) {

	type FormData struct {
		Urls []string `json:"data"`
	}

	data := &FormData{}
	if err := c.ShouldBindJSON(data); err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	prxs, err := ctrl.proxyService.AddRange(data.Urls)

	until.WriteResponse(c, 200, gin.H{
		"count":   len(prxs),
		"proxies": prxs,
	}, err)
}

// freeHandler
// Получаем свободное прокси
func (ctrl *ProxyController) freeHandler(c *gin.Context) {
	proxy, err := ctrl.proxyService.Free()
	until.WriteResponse(c, 200, gin.H{
		"proxy": proxy,
	}, err)
}

/*
________________________

cleatAllProxy
Очищаем все прокси в базе


_______________________
*/
func (ctrl *ProxyController) clearAllHandler(c *gin.Context) {

	err := ctrl.proxyService.Clear()

	ok := true
	if err != nil {
		ok = false
	}

	until.WriteResponse(c, 200, gin.H{
		"ok": ok,
	}, err)
}
