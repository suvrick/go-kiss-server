package client

import (
	"github.com/suvrick/go-kiss-server/game/packets/encode"
	"github.com/suvrick/go-kiss-server/game/parser"
)

// "LBBS,BSIIBSBSBS",
// LOGIN(4); net_id:L, SocType:B, device:B, auth_key:S, oauth:B, session_key:S,
// referrer:I, tag:I, appicationID:B, timestamp:S, language:B, utm_source:S, sex:B, captcha:S

// LoginClientPacket ...
type LoginClientPacket struct {
	Type        int16
	DeviceType  byte
	ID          int64
	SocialType  byte
	DeviceType2 byte
	Token       string
	IOAth       byte
	Token2      string
	Referer     int
	Tag         int
	AppID       byte
	Timestamp   string
	Language    byte
	UtmData     string
	Gender      byte
	Captcha     string
}

// NewLoginClientPacket ...
func NewLoginClientPacket(login *parser.LoginParams) encode.ClientPacket {

	var isOAuth byte
	if len(login.Token2) > 0 {
		isOAuth = 1
	}

	return &LoginClientPacket{
		Type:        4,
		DeviceType:  0x04,
		ID:          login.LoginID,
		SocialType:  uint8(login.SocialCode),
		DeviceType2: 0x04,
		Token:       login.Token,
		IOAth:       isOAuth,
		Token2:      login.Token2,
		Referer:     0,
		Tag:         13,
		AppID:       0,
		Timestamp:   "",
		Language:    0x01,
		UtmData:     "",
		Gender:      0,
		Captcha:     "",
	}
}

// Bytes ...
func (p *LoginClientPacket) Bytes() []byte {
	return encode.Load(p)
}
