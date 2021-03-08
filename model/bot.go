package model

import (
	"github.com/suvrick/go-kiss-core/parser"
)

// Bot ...
type Bot struct {
	ID      int
	UserID  string
	Result  string
	DateUse string
	parser.LoginParams
}

// NewBot ...
func NewBot(lp parser.LoginParams) *Bot {
	return &Bot{
		LoginParams: lp,
	}
}
