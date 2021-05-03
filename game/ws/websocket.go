package ws

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/packets/decode"
	"github.com/suvrick/go-kiss-server/game/packets/encode"

	wss "github.com/gorilla/websocket"
)

// GameSock ...
type GameSock struct {
	client         *wss.Conn
	msgID          uint32
	bot            *models.Bot
	proxy          *models.Proxy
	botChanUpdater chan *models.Bot
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
	}

	return gs
}

// Go start game
func (gs *GameSock) Go() {

	if gs.bot.IsError {
		gs.close()
		return
	}

	gs.bot.LogINFO("Go", "Try connection")
	err := gs.connect()

	if err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("Go", err.Error())
		gs.error("Go", err)
		return
	}

	gs.loginSend()
	gs.readMessage()
}

func (gs *GameSock) connect() error {

	dialer := wss.Dialer{}

	//if gs.bot.Proxy != nil {
	dialer = wss.Dialer{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "zproxy.lum-superproxy.io:22222",
			User:   url.UserPassword("lum-customer-c_07f044e7-zone-static", "dodwwsy0fhb00"),
		}),
		HandshakeTimeout: (time.Second * 60),
	}
	//}

	con, _, err := dialer.Dial(host, nil)

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

		log.Printf("Recv >> msgType: %d,msgID: %d, msgLen: %d", msgType, msgID, msgLen)

		switch msgType {
		case 4:
			ok := gs.loginReceive(reader)
			if !ok {
				gs.error("readMessage", errors.New("Auth error"))
				return
			}
		case 5:
			gs.infoReceive(reader)
			gs.close()
			return
		case 7:
			gs.balanceReceive(reader)
		case 13:
			gs.gameListRewardsReceive(reader)
		case 17:
			gs.bonusReceive(reader)
		}
	}
}

func (gs *GameSock) sendMessage(pack encode.ClientPacket) {

	msg := pack.Bytes()

	IDArr := make([]byte, 0)
	IDArr = encode.WriteNumber(IDArr, uint64(gs.msgID))
	lengthMsg := len(msg) + len(IDArr)

	data := make([]byte, 0)
	data = append(data, encode.WriteNumber(data, uint64(lengthMsg))...)
	data = append(data, IDArr...)
	data = append(data, msg...)

	log.Printf("%v", data)

	err := gs.client.WriteMessage(wss.BinaryMessage, data)
	gs.msgID++

	if err != nil {
		gs.error("sendMessage", err)
		return
	}
}

func (gs *GameSock) error(funcName string, err error) {
	//gs.bot.Error = fmt.Sprintf("[%s] >> %s", funcName, err.Error())
	gs.bot.IsError = true
	gs.close()
}

func (gs *GameSock) close() {

	if gs.client != nil {
		gs.client.Close()
		gs.client = nil
	}
}
