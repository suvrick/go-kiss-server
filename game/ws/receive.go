package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/suvrick/go-kiss-server/game/packets/server"
)

//Авторизация игрока
// type 4
func (gs *GameSock) loginReceive(msg io.Reader) bool {

	gs.bot.LogINFO("LoginServerPacket", "Try parse server packet")

	login := server.LoginServerPacket{}

	err := login.Parse(msg)

	gs.bot.Result = login.Result.ToString()
	gs.bot.LogINFO("LoginServerPacket", fmt.Sprintf("Set login status: %v", gs.bot.Result))

	if err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("LoginServerPacket", "Error parse packet")
		return false
	}

	if server.LOGIN_SUCCESS == login.Result {

		gs.bot.InnerID = login.InnerID
		gs.bot.Balance = login.Balance

		gs.bot.LogINFO("LoginServerPacket", fmt.Sprintf("Set innerID: %v", login.InnerID))
		gs.bot.LogINFO("LoginServerPacket", fmt.Sprintf("Set balance: %v", login.Balance))
		return true
	}

	return false
}

//Информация о игроке (парсится по 2 маскам!)
//type 5
func (gs *GameSock) infoReceive(msg io.Reader) {
	gs.bot.LogINFO("InfoServerPacket", "Try parse server packet")
	info := server.InfoServerPacket{}

	if err := info.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("InfoServerPacket", "Error parse server packet")
		return
	}

	gs.bot.Name = info.Name
	gs.bot.Avatar = strings.Split(info.Avatar, "?")[0]
	gs.bot.Profile = info.Profile

	gs.bot.LogINFO("InfoServerPacket", fmt.Sprintf("Set name: %v", gs.bot.Name))
	gs.bot.LogINFO("InfoServerPacket", fmt.Sprintf("Set avatar: %v", gs.bot.Avatar))
	gs.bot.LogINFO("InfoServerPacket", fmt.Sprintf("Set profile: %v", gs.bot.Profile))

}

//Пакет с массиво наград игроку
// type = 13
func (gs *GameSock) gameListRewardsReceive(msg io.Reader) int32 {

	gs.bot.LogINFO("GameRewardsListServerPacket", "Try parse server packet")

	//Запрос на взятия ежедневного бонуса должен быть обработан до этого пакета
	rewardsList := server.GameRewardsListServerPacket{}
	if err := rewardsList.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("GameRewardsListServerPacket", "Error parse packet")
	}

	type Content struct {
		Coints int `json:"coins"`
	}

	type Captions struct {
		Ru string `json:"ru"`
		En string `json:"en"`
	}

	type Result struct {
		ID       int      `json:"id"`
		Content  Content  `json:"content"`
		Captions Captions `json:"captions"`
		Type     string   `json:"type"`
	}

	result := &Result{}

	if rewardsList.ID != 0 {
		if err := json.Unmarshal([]byte(rewardsList.Json), result); err != nil {
			gs.bot.LogERROR("GameRewardsListServerPacket", err.Error())
		}
	}

	//fmt.Println(rewardsList.Json)

	gs.bot.LogINFO("GameRewardsListServerPacket", fmt.Sprintf("id: %v, count: %v, captions: %v", result.ID, result.Content, result.Captions.Ru))
	return rewardsList.ID
}

//Пакет ежедневного бонуса
//type 17
func (gs *GameSock) bonusReceive(msg io.Reader) {

	gs.bot.LogINFO("BonusServerPacket", "Try parse server packet")

	bonus := server.BonusServerPacket{}
	if err := bonus.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("BonusServerPacket", "Error parse server packet")
		return
	}

	if bonus.CanCollect == 0x01 {
		gs.bot.IsBonus = true
	}

	gs.bot.BonusDay = int32(bonus.Day)

	gs.bot.LogINFO("BonusServerPacket", fmt.Sprintf("Set can bonus: %v", gs.bot.IsBonus))
	gs.bot.LogINFO("BonusServerPacket", fmt.Sprintf("Set bonus day: %d", gs.bot.BonusDay))
}

//Пакет обновления баланса игрока
//type 7
func (gs *GameSock) balanceReceive(msg io.Reader) {

	gs.bot.LogINFO("BalanceServerPacket", "Try parse server packet")

	//fmt.Println(msg)

	balance := server.BalanceServerPacket{}

	if err := balance.Parse(msg); err != nil {
		gs.bot.IsError = true
		gs.bot.LogERROR("BalanceServerPacket", "Error parse packet")
		return
	}

	gs.bot.Balance = balance.Bottles
	gs.bot.LogINFO("BalanceServerPacket", fmt.Sprintf("Update balance: %d", gs.bot.Balance))
}
