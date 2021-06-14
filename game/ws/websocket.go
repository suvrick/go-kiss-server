package ws

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/packets/decode"
	"github.com/suvrick/go-kiss-server/game/packets/encode"

	wss "github.com/gorilla/websocket"
)

type IProxyManager interface {
	GetString() (string, error)
	UpdateString(host string, isError bool) error
}

// GameSock ...
type GameSock struct {
	client         *wss.Conn
	msgID          uint32
	bot            *models.Bot
	botChanUpdater chan *models.Bot

	packets      []*encode.ClientPacket
	proxyManager IProxyManager
	debug        bool

	proxyHost     string
	proxyLogin    string
	proxyPassword string

	end bool
}

const host = "wss://bottlews.itsrealgames.com"

// NewSocket созадает и заполняет новую структуру
// p - ссылка на Player
// updater - канал типа Player
func NewSocket(b *models.Bot) *GameSock {

	gs := &GameSock{
		client: nil,
		msgID:  0,
		bot:    b,
		debug:  false,
	}

	return gs
}

// NewSocket созадает и заполняет новую структуру
// p - ссылка на Player
// updater - канал типа Player
func NewSocketWithAdditionPacket(b *models.Bot, p []*encode.ClientPacket) *GameSock {

	gs := &GameSock{
		client:  nil,
		msgID:   0,
		bot:     b,
		debug:   true,
		packets: p,
	}

	return gs
}

func (gs *GameSock) SetProxyManager(p IProxyManager) {
	gs.proxyManager = p
}

// Go start game
func (gs *GameSock) Go() {

	gs.bot.Logger = make([]models.LoggerLine, 0)

	gs.bot.LogINFO("Go", "Try connection")
	err := gs.connect()

	if err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("Go", err.Error())
		gs.error("Go", err)
		return
	}

	go func(g *GameSock) {
		time.Sleep(time.Minute * time.Duration(5))

		if g == nil {
			return
		}

		g.bot.LogERROR("Go", "close socket by timeout")
		g.close()
	}(gs)

	gs.loginSend()
	gs.readMessage()
}

func (gs *GameSock) connect() error {

	var dialer wss.Dialer
	var array []string

	if gs.proxyManager != nil && gs.packets == nil {

		proxy, err := gs.proxyManager.GetString()

		if err != nil {
			return err
		}

		array = strings.Split(proxy, ":")

		if len(array) != 4 {
			return errors.New("invalid proxy")
		}

		dialer = wss.Dialer{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%s", array[0], array[1]),
				User:   url.UserPassword(array[2], array[3]),
			}),
			HandshakeTimeout: (time.Second * 60),
		}

		gs.bot.LogINFO("connect", fmt.Sprintf("Set proxy %s", array[2]))
	} else {
		dialer = wss.Dialer{}
	}

	con, _, err := dialer.Dial(host, nil)

	if err != nil {
		gs.proxyManager.UpdateString(array[2], true)
		return gs.connect()
	}

	//gs.proxyManager.UpdateString(array[2], false)
	gs.client = con
	return err
}

func (gs *GameSock) readMessage() {

	for {

		if gs.client == nil {
			gs.error("readMessage", errors.New("Connection is nil"))
			return
		}

		_, msg, err := gs.client.ReadMessage()

		if err != nil {
			gs.error("readMessage", err)
			return
		}

		if len(msg) < 3 {
			continue
		}

		reader := bytes.NewReader(msg)

		msgLen, _ := decode.ReadVarUint(reader, 32)
		msgID, _ := decode.ReadVarUint(reader, 32)
		msgType, _ := decode.ReadVarUint(reader, 16)

		if gs.debug {
			log.Printf("Recv >> msgType: %d,msgID: %d, msgLen: %d\n", msgType, msgID, msgLen)
		}

		//fmt.Printf("Recv >> msgType: %d,msgID: %d, msgLen: %d\n", msgType, msgID, msgLen)

		switch msgType {
		case 4:
			ok := gs.loginReceive(reader)
			if !ok {
				gs.error("readMessage", errors.New("Not authenticated"))
				return
			}

			gs.bonusSend()

			if len(gs.packets) > 0 {
				gs.additionPacketSend()
			}

		case 5:
			gs.infoReceive(reader)
		case 7:
			gs.balanceReceive(reader)
		case 13:
			id := gs.gameListRewardsReceive(reader)
			if id == 0 {
				continue
			}

			gs.getRewardSend(id)
		case 9:
			go func() {
				time.Sleep(time.Second * time.Duration(5))
				gs.close()
			}()
		case 17:
			gs.bonusReceive(reader)
		}

	}
}

func (gs *GameSock) sendMessage(pack encode.ClientPacket) {

	msg := pack.Bytes()

	msgID_array := make([]byte, 0)
	msgID_array = encode.WriteNumber(msgID_array, uint64(gs.msgID))
	lengthMsg := len(msg) + len(msgID_array)

	data := make([]byte, 0)
	data = append(data, encode.WriteNumber(data, uint64(lengthMsg))...)
	data = append(data, msgID_array...)
	data = append(data, msg...)

	if gs.debug {
		log.Printf("Send >> id: %v, data: %v\n", gs.msgID, data)
	}

	err := gs.client.WriteMessage(wss.BinaryMessage, data)
	gs.msgID++

	if err != nil {
		gs.error("sendMessage", err)
		return
	}
}

func (gs *GameSock) error(funcName string, err error) {
	gs.bot.LogERROR(funcName, err.Error())
	gs.bot.IsError = true
	gs.close()
}

func (gs *GameSock) close() {

	if gs.client != nil {
		gs.client.Close()
		gs.client = nil
	}

	if gs.bot.IsError {
		gs.bot.LogERROR("close", "error close connection")
	} else {
		gs.bot.LogINFO("close", "normal close connection")
	}
}
