package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// ViewClientPacket ...
type ViewClientPacket struct {
	Type       int16
	DeviceType byte
	TargetID   int
}

// NewViewClientPacket ...
func NewViewClientPacket(targetID int) encode.ClientPacket {
	return &ViewClientPacket{
		Type:       17,
		DeviceType: 0x04,
		TargetID:   targetID,
	}
}

// Bytes ...
func (p *ViewClientPacket) Bytes() []byte {
	return encode.Load(p)
}
