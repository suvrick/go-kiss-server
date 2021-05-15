package server

import (
	"io"

	"github.com/suvrick/go-kiss-server/game/packets/decode"
)

type BonusServerPacket struct {
	CanCollect byte
	Day        byte
}

func (pack *BonusServerPacket) Parse(buffer io.Reader) error {
	return decode.Fill(pack, buffer)
}
