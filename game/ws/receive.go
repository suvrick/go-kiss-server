package ws

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/suvrick/go-kiss-server/game/packets/server"
)

//Авторизация игрока
// type 4
func (gs *GameSock) loginReceive(msg io.Reader) bool {

	login := server.LoginServerPacket{}
	if err := login.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("LoginServerPacket", "Error parse packet")
		return false
	}

	switch login.Result {
	case server.LOGIN_SUCCESS:

		gs.bot.Result = "SUCCESS"
		gs.bot.InnerID = login.InnerID
		gs.bot.Balance = login.Balance

		//send client packet bonus day
		gs.bonusSend()

		//send client packet get reward(25)
		// id 25 === bonus day
		gs.getRewardSend(25)

		return true
	case server.LOGIN_EXIST:
		gs.bot.Result = "EXIST"
	case server.LOGIN_BLOCKED:
		gs.bot.Result = "BLOCKED"
	case server.LOGIN_FAILED:
		gs.bot.Result = "FAILED"
	case server.LOGIN_CAPTCHA:
		gs.bot.Result = "CAPTCHA"
	default:
		gs.bot.Result = "XZ"
		break
	}

	return false
}

//Информация о игроке (парсится по 2 маскам!)
//type 5
func (gs *GameSock) infoReceive(msg io.Reader) {

	info := server.InfoServerPacket{}

	if err := info.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("InfoServerPacket", "Error parse packet")
		return
	}

	gs.bot.Name = info.Name
	gs.bot.Avatar = strings.Split(info.Avatar, "?")[0]
	gs.bot.Profile = info.Profile

}

//Пакет с массиво наград игроку
// type = 13
func (gs *GameSock) gameListRewardsReceive(msg io.Reader) {
	//Запрос на взятия ежедневного бонуса должен быть обработан до этого пакета
	rewardsList := server.GameRewardsListServerPacket{}
	if err := rewardsList.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("GameRewardsListServerPacket", "Error parse packet")
		return
	}

	//fmt.Printf("%v\n", *rewardsList)
	type Content struct {
		Coint int
	}

	type Captions struct {
		Ru string
		En string
	}

	type Result struct {
		ID       int
		Content  Content
		Captions Captions
		Type     string
	}

	result := &Result{}

	if rewardsList.ID != 0 {
		if err := json.Unmarshal([]byte(rewardsList.Json), result); err == nil {
			log.Println(result)
			gs.getRewardSend(rewardsList.ID)
		}
	}
}

//Пакет ежедневного бонуса
//type 17
func (gs *GameSock) bonusReceive(msg io.Reader) {

	bonus := server.BonusServerPacket{}
	if err := bonus.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("BonusServerPacket", "Error parse packet")
		return
	}

	if bonus.CanCollect == 0x01 {
		gs.bot.IsBonus = true
		day := int32(bonus.Day)

		if day == 7 {
			gs.bot.BonusDay = 15
		} else {
			gs.bot.BonusDay = day + 1
		}
	}
}

//Пакет обновления баланса игрока
//type 7
func (gs *GameSock) balanceReceive(msg io.Reader) {
	balance := server.BalanceServerPacket{}
	if err := balance.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("BalanceServerPacket", "Error parse packet")
		return
	}

	gs.bot.Balance = balance.Bottles
}
