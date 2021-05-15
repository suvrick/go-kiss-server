package client

import "github.com/suvrick/go-kiss-server/game/packets/encode"

// "IIIIB,SS",
// BUY(6); good_id:I, cost:I, target_id:I, data:I, price_type:B, hash:S, params: S

// PrizeClientPacket ...
type PrizeClientPacket struct {
	Type       int16
	DeviceType byte
	Good_id    int
	Cost       int
	Target_id  int
	Data       int
	Price_type byte
	Hash       string
	Params     string
}

// NewPrizeClientPacket ...
func NewPrizeClientPacket() encode.ClientPacket {
	return &PrizeClientPacket{
		Type:       6,
		DeviceType: 0x04,
		Good_id:    0,
		Cost:       0,
		Target_id:  0,
		Data:       0,
		Price_type: 0,
		Hash:       "",
		Params:     "",
	}
}

// Bytes ...
func (p *PrizeClientPacket) Bytes() []byte {
	return encode.Load(p)
}
