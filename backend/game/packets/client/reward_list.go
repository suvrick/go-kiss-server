package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// RewardListClientPacket ...
type RewardListClientPacket struct {
	Type       int16
	DeviceType byte
}

// NewRewardListClientPacket ...
func NewRewardListClientPacket() encode.ClientPacket {
	return &RewardListClientPacket{
		Type:       163,
		DeviceType: 0x04,
	}
}

// Bytes ...
func (p *RewardListClientPacket) Bytes() []byte {
	return encode.Load(p)
}
