package server

import (
	"io"

	"github.com/suvrick/go-kiss-server/game/packets/decode"
)

type SERVER_AUTH_RESULT byte

const (
	LOGIN_SUCCESS       SERVER_AUTH_RESULT = 0x00
	LOGIN_FAILED        SERVER_AUTH_RESULT = 0x01
	LOGIN_EXIST         SERVER_AUTH_RESULT = 0x02
	LOGIN_BLOCKED       SERVER_AUTH_RESULT = 0x03
	LOGIN_WRONG_VERSION SERVER_AUTH_RESULT = 0x04
	LOGIN_NO_SEX        SERVER_AUTH_RESULT = 0x05
	LOGIN_CAPTCHA       SERVER_AUTH_RESULT = 0x06

	LOGIN_WAIT_AUTHORIZATION SERVER_AUTH_RESULT = 0xfc
	LOGIN_ERROR              SERVER_AUTH_RESULT = 0xff
)

func (s SERVER_AUTH_RESULT) ToString() string {
	switch s {
	case LOGIN_SUCCESS:
		return "SUCCESS"
	case LOGIN_EXIST:
		return "EXIST"
	case LOGIN_BLOCKED:
		return "BLOCKED"
	case LOGIN_FAILED:
		return "FAILED"
	case LOGIN_CAPTCHA:
		return "CAPTCHA"
	default:
		return "XZ"
	}
}

// "B,IIBI[B]IIIISBIBS"
//LOGIN(4); status:B, inner_id:I, balance:I, invited:B, logout_time:I,
//[flags:B], games_count:I, kisses_daily_count:I, last_payment_time:I,
//subscribe_expires:I, params:S, sex_set:B, server_time:I, first_login:B, photos_hash:S

//LoginServerPacket ...
type LoginServerPacket struct {
	Result  SERVER_AUTH_RESULT
	InnerID int32
	Balance int32
}

// Parse ...
func (pack *LoginServerPacket) Parse(buffer io.Reader) error {
	return decode.Fill(pack, buffer)
}
