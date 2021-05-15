package server

import (
	"io"

	"github.com/suvrick/go-kiss-server/game/packets/decode"
)

// BALANCE(7); bottles:I, reason:B

type BalanceServerPacket struct {
	Bottles int32
	//Reason  byte
}

func (pack *BalanceServerPacket) Parse(buffer io.Reader) error {
	return decode.Fill(pack, buffer)
}
