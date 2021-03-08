package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
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
	proxy := ctrl.router.Group("/proxy")
	proxy.POST("/add", ctrl.addHandler)
	proxy.GET("/all", ctrl.allHandler)
	proxy.GET("/free", ctrl.freeHandler)
	proxy.GET("/deleteall", ctrl.deleteAllHandler)
	return ctrl
}

func (ctrl *ProxyController) allHandler(c *gin.Context) {
	prxs, err := ctrl.proxyService.All()

	until.WriteResponse(c, 200, gin.H{
		"count":   len(prxs),
		"proxies": prxs,
	}, err)
}

func (ctrl *ProxyController) addHandler(c *gin.Context) {

	type U struct {
		Urls []string `json:"urls"`
	}

	urls := &U{}
	if err := c.ShouldBindJSON(urls); err != nil {
		until.WriteResponse(c, 200, nil, errors.ErrInvalidParam)
		return
	}

	prxs, err := ctrl.proxyService.AddRange(urls.Urls)

	until.WriteResponse(c, 200, gin.H{
		"count":   len(prxs),
		"proxies": prxs,
	}, err)
}

func (ctrl *ProxyController) freeHandler(c *gin.Context) {
	proxy, err := ctrl.proxyService.Free()
	until.WriteResponse(c, 200, gin.H{
		"proxy": proxy,
	}, err)
}

func (ctrl *ProxyController) deleteAllHandler(c *gin.Context) {
	err := ctrl.proxyService.DeleteAll()
	ok := true
	if err != nil {
		ok = false
	}

	until.WriteResponse(c, 200, gin.H{
		"ok": ok,
	}, err)
}
