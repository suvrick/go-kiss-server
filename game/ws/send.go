package ws

import (
	"github.com/suvrick/go-kiss-server/game/packets/client"
)

func (gs *GameSock) loginSend() {
	data := client.NewLoginClientPacket(&gs.bot.LoginParams)
	gs.sendMessage(data)
}

func (gs *GameSock) bonusSend() {
	//log.Println("send -> Bonus")
	data := client.NewBonusClientPacket()
	gs.sendMessage(data)
}

func (gs *GameSock) rewardListSend() {
	//log.Println("send -> Get rewards list")
	data := client.NewRewardListClientPacket()
	gs.sendMessage(data)
}

func (gs *GameSock) getRewardSend(id int32) {
	//log.Println("send -> Get reward by id: ", id)
	data := client.NewRewardClientPacket(id)
	gs.sendMessage(data)
}

func (gs *GameSock) postAwardsSend() {
	//log.Println("send -> Post award")
	data := client.NewPostAwardsClientPacket()
	gs.sendMessage(data)
}
