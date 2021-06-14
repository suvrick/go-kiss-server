package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/suvrick/go-kiss-server/game/packets/client"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/session"
)

type WsController struct {
	router      *gin.Engine
	userService *services.UserService
	botService  *services.BotService
}

type WSConn struct {
	ws         *websocket.Conn
	user       *model.User
	msg        chan []byte
	botService *services.BotService

	locker *sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWsController(r *gin.Engine, u_service *services.UserService, b_service *services.BotService) *WsController {
	ctrl := &WsController{
		router:      r,
		userService: u_service,
		botService:  b_service,
	}

	ctrl.router.Use(middlewares.AuthMiddleware())
	ctrl.router.GET("/ws", ctrl.acceptHandler)

	return ctrl
}

func (ctrl *WsController) acceptHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}

	wsConn := WSConn{
		ws:         conn,
		user:       session.GetUser(c),
		msg:        make(chan []byte),
		botService: ctrl.botService,

		locker: &sync.Mutex{},
	}

	go wsConn.reader(ctrl.botService)
}

func (c *WSConn) reader(botService *services.BotService) {
	defer func() {
		fmt.Println("close socket")
		c.ws.Close()
	}()

	c.cmdSelf()

	for {

		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		p := &ClientPacket{}
		err = json.Unmarshal(msg, p)

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		//fmt.Println(p)

		switch p.Type {
		case ALL_BOT_SEND:
			c.cmdAllBot()
			break
		case ADD_BOT_SEND:
			c.cmdAddTask()
			go c.cmdAddBot(p.Data)
			break
		case REMOVE_BOT_SEND:
			c.cmdAddTask()
			go c.cmdRemoveBot(p.Data)
			break
		case UPDATE_BOT_SEND:
			c.cmdAddTask()
			go c.cmdUpdateBot(p.Data)
			break
		case PRIZE_BOT_SEND:
			go c.cmdPrizeBot(p.Data)
			break
		case VIEW_BOT_SEND:
			go c.cmdViewBot(p.Data)
			break
		}

	}
}

func (c *WSConn) cmdSelf() {
	c.send(SELF_RECV, map[string]interface{}{
		"user": c.user,
	})
}

func (c *WSConn) cmdAddTask() {
	c.send(ADD_TASK_RECV, nil)
}

func (c *WSConn) cmdRemoveTask() {
	c.send(REMOVE_TASK_RECV, nil)
}

func (c *WSConn) cmdAllBot() {
	bots, err := c.botService.All(c.user)

	if err != nil {
		c.send(ERROR_RECV, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.send(ALL_BOT_RECV, map[string]interface{}{
		"bots": bots,
	})
}

func (c *WSConn) cmdAddBot(data map[string]interface{}) {

	url := data["url"].(string)

	bot, err := c.botService.Add(c.user, url)

	if err != nil {
		c.send(ERROR_RECV, map[string]interface{}{
			"error": err.Error(),
		})
		c.cmdRemoveTask()
		return
	}

	c.send(ADD_BOT_RECV, map[string]interface{}{
		"bot": bot,
	})

	c.cmdRemoveTask()
}

func (c *WSConn) cmdRemoveBot(data map[string]interface{}) {
	uid := data["uid"].(string)
	err := c.botService.Delete(uid, c.user)
	if err != nil {
		c.send(ERROR_RECV, map[string]interface{}{
			"error": err.Error(),
		})
		c.cmdRemoveTask()
		return
	}

	c.send(REMOVE_BOT_RECV, map[string]interface{}{
		"uid": uid,
	})

	c.cmdRemoveTask()
}

func (c *WSConn) cmdUpdateBot(data map[string]interface{}) {

	uid := data["uid"].(string)

	bot, err := c.botService.UpdateByID(uid, c.user)

	if err != nil {
		c.send(ERROR_RECV, map[string]interface{}{
			"error": err.Error(),
		})
		c.cmdRemoveTask()
		return
	}

	c.send(UPDATE_BOT_RECV, map[string]interface{}{
		"bot": bot,
	})

	c.cmdRemoveTask()
}

/*

	ПОДАРКИ

*/
func (c *WSConn) cmdPrizeBot(data map[string]interface{}) {

	uids := data["uids"].([]interface{})

	id, ok := data["target_id"].(float64)
	if !ok {
		return
	}

	count, ok := data["count"].(float64)
	if !ok {
		return
	}

	prize := client.NewPrizeClientPacket(
		int(data["good_id"].(float64)),
		int(data["cost"].(float64)),
		int(id),
		int(data["data"].(float64)),
		byte(data["price_type"].(float64)),
		int(count),
		data["hash"].(string),
		data["params"].(string),
	)

	for _, v := range uids {

		go func(v interface{}, c *WSConn) {
			c.cmdAddTask()

			bot, err := c.botService.SendPrize2(v.(string), c.user, &prize)

			if err != nil {
				c.send(ERROR_RECV, map[string]interface{}{
					"error": err.Error(),
				})
				c.cmdRemoveTask()
				return
			}

			c.send(UPDATE_BOT_RECV, map[string]interface{}{
				"bot": bot,
			})

			c.cmdRemoveTask()
		}(v, c)
	}

}

/*

	ПРОСМОТРЫ

*/
func (c *WSConn) cmdViewBot(data map[string]interface{}) {

	uids := data["uids"].([]interface{})
	targetID, ok := data["target_id"].(float64)
	if !ok {
		return
	}

	pack := client.NewViewClientPacket(int(targetID))

	for _, v := range uids {

		go func(v interface{}, c *WSConn) {
			c.cmdAddTask()

			bot, err := c.botService.SendPrize2(v.(string), c.user, &pack)

			if err != nil {
				c.send(ERROR_RECV, map[string]interface{}{
					"error": err.Error(),
				})
				c.cmdRemoveTask()
				return
			}

			c.send(UPDATE_BOT_RECV, map[string]interface{}{
				"bot": bot,
			})

			c.cmdRemoveTask()
		}(v, c)
	}
}

func (c *WSConn) send(t PacketServerType, d map[string]interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()

	msg := createServerPacket(t, d)
	c.ws.WriteMessage(websocket.TextMessage, msg)
}

func createServerPacket(t PacketServerType, d map[string]interface{}) []byte {

	p := ServerPacket{
		Type: t,
		Data: d,
	}

	data, err := json.Marshal(p)

	if err != nil {
		return nil
	}

	return data
}

type ServerPacket struct {
	Type PacketServerType       `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type ClientPacket struct {
	Type PacketClientType       `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type PacketServerType string

const (
	SELF_RECV PacketServerType = "SELF_RECV"

	ADD_BOT_RECV    PacketServerType = "ADD_BOT_RECV"
	ALL_BOT_RECV    PacketServerType = "ALL_BOT_RECV"
	REMOVE_BOT_RECV PacketServerType = "REMOVE_BOT_RECV"
	UPDATE_BOT_RECV PacketServerType = "UPDATE_BOT_RECV"

	ADD_TASK_RECV    PacketServerType = "ADD_TASK_RECV"
	REMOVE_TASK_RECV PacketServerType = "REMOVE_TASK_RECV"

	ERROR_RECV PacketServerType = "ERROR_RECV"
)

type PacketClientType string

const (
	ADD_BOT_SEND    PacketClientType = "ADD_BOT_SEND"
	ALL_BOT_SEND    PacketClientType = "ALL_BOT_SEND"
	REMOVE_BOT_SEND PacketClientType = "REMOVE_BOT_SEND"
	UPDATE_BOT_SEND PacketClientType = "UPDATE_BOT_SEND"

	PRIZE_BOT_SEND PacketClientType = "PRIZE_BOT_SEND"
	VIEW_BOT_SEND  PacketClientType = "VIEW_BOT_SEND"
)
