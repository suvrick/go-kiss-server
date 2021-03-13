package controllers

import (
	"encoding/binary"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/services"
)

type AutoKissController struct {
	router               *gin.Engine
	kissUserService      *services.AutoKissService
	stateDownloadService *services.StateDownloadService
}

func NewAutoKissController(r *gin.Engine, kus *services.AutoKissService, sds *services.StateDownloadService) *AutoKissController {
	ctrl := &AutoKissController{
		router:               r,
		kissUserService:      kus,
		stateDownloadService: sds,
	}

	kiss := ctrl.router.Group("autokiss")
	kiss.GET("", ctrl.indexHandler)
	kiss.GET("/zip", ctrl.downloadZipHandler)
	kiss.GET("/js", ctrl.jsHandler)
	kiss.GET("/init/:userID", ctrl.initKissUserHandler)
	kiss.POST("/who/:userID", ctrl.whoHandler)

	return ctrl
}

func (ctrl *AutoKissController) whoHandler(context *gin.Context) {

	id := context.Param("userID")

	userID, _ := strconv.Atoi(id)

	body := context.Request.Body
	data, err := ioutil.ReadAll(body)

	if err != nil {
		context.JSON(http.StatusOK, &model.KissResponse{Code: 0})
		return
	}

	if len(data) < 6 {
		context.JSON(http.StatusOK, &model.KissResponse{Code: 0})
		return
	}

	pType := binary.LittleEndian.Uint16(data[4:6])

	res := ctrl.kissUserService.Do(userID, pType, data)

	context.JSON(http.StatusOK, res)
}

func (ctrl *AutoKissController) indexHandler(c *gin.Context) {
	c.File("www/autokiss/index.html")
}

func (ctrl *AutoKissController) jsHandler(c *gin.Context) {

	file, _ := ioutil.ReadFile("www/autokiss/in.js")
	file2, _ := ioutil.ReadFile("www/autokiss/style.css")

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
		"js":     string(file),
		"css":    string(file2),
	})
}

func (ctrl *AutoKissController) downloadZipHandler(c *gin.Context) {

	ip := c.ClientIP()

	state, _ := ctrl.stateDownloadService.FindDownloadState(ip)

	if state.IP == "" {
		state.IP = ip
		ctrl.stateDownloadService.AddDownloadState(ip)
	}

	state.Count++
	state.Date = time.Now().Format(time.RFC822)

	ctrl.stateDownloadService.UpdateDownloadState(state)

	c.File("www/autokiss/autokiss.zip")
}

func (ctrl *AutoKissController) initKissUserHandler(c *gin.Context) {

	id := c.Param("userID")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, gin.H{"result": "fail"})
	}

	user, err := ctrl.kissUserService.FindByIDKissUser(userID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ctrl.kissUserService.AddKissUser(userID)
		}
	}

	if user.UserID == 0 {
		user = &model.KissUser{
			UserID:  userID,
			IsTrial: true,
			DateUse: time.Now().Format(time.RFC822),
		}
	}

	ctrl.kissUserService.UpdateKissUser(user)

	c.JSON(200, gin.H{"result": "ok", "user": user})
}
