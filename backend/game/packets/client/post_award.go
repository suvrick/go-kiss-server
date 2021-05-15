package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// PostAwardsClientPacket ...
type PostAwardsClientPacket struct {
	Type       uint32
	DeviceType byte
	TypeByte   byte
}

// NewPostAwardsClientPacket ...
func NewPostAwardsClientPacket() encode.ClientPacket {
	return &PostAwardsClientPacket{
		Type:       163,
		DeviceType: 0x04,
		TypeByte:   0x01,
	}
}

// Bytes ...
func (p *PostAwardsClientPacket) Bytes() []byte {
	return encode.Load(p)
}
