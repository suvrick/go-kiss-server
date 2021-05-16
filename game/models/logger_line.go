package models

import "gorm.io/gorm"

type LoggerLine struct {
	gorm.Model
	BotID uint
	Line  string
}
