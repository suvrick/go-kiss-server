package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// RewardClientPacket ...
type RewardClientPacket struct {
	Type       int16
	DeviceType byte
	RewardID   int
}

// NewRewardClientPacket ...
func NewRewardClientPacket(rewardID int32) encode.ClientPacket {
	return &RewardClientPacket{
		Type:       11,
		DeviceType: 0x04,
		RewardID:   int(rewardID),
	}
}

// Bytes ...
func (p *RewardClientPacket) Bytes() []byte {
	return encode.Load(p)
}
