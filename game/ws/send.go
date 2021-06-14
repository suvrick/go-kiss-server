package ws

import (
	"fmt"

	"github.com/suvrick/go-kiss-server/game/packets/client"
)

func (gs *GameSock) loginSend() {
	data := client.NewLoginClientPacket(&gs.bot.LoginParams)
	gs.bot.LogINFO("loginSend", "Try send login")
	gs.sendMessage(data)
}

func (gs *GameSock) bonusSend() {
	gs.bot.LogINFO("bonusSend", "Try send daily bonus")
	data := client.NewBonusClientPacket()
	gs.sendMessage(data)
}

func (gs *GameSock) additionPacketSend() {

	gs.bot.LogINFO("additionPacketSend", "Try send addiction packet send")
	for _, p := range gs.packets {
		gs.sendMessage(*p)
	}
}

func (gs *GameSock) rewardListSend() {
	gs.bot.LogINFO("rewardListSend", "Try send reward list")
	data := client.NewRewardListClientPacket()
	gs.sendMessage(data)
}

func (gs *GameSock) getRewardSend(id int32) {
	gs.bot.LogINFO("getRewardSend", fmt.Sprintf("Try send reward get by id: %v", id))
	data := client.NewRewardClientPacket(id)
	gs.sendMessage(data)
}

func (gs *GameSock) postAwardsSend() {
	gs.bot.LogINFO("postAwardsSend", "Try send post award")
	data := client.NewPostAwardsClientPacket()
	gs.sendMessage(data)
}

func (gs *GameSock) closeSend() {
	gs.bot.LogINFO("closeSend", "send packet by close socket")
	data := client.NewCloseClientPacket()
	gs.sendMessage(data)
}
