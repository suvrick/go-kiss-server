package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// CloseClientPacket ...
type CloseClientPacket struct {
	Type       int16
	DeviceType byte
}

// NewBonusClientPacket ...
func NewCloseClientPacket() encode.ClientPacket {
	return &CloseClientPacket{
		Type:       100,
		DeviceType: 0x04,
	}
}

// Bytes ...
func (p *CloseClientPacket) Bytes() []byte {
	return encode.Load(p)
}
