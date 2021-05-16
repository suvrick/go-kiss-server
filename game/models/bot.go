package models

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-server/game/parser"
	"github.com/suvrick/go-kiss-server/until"
	"gorm.io/gorm"
)

type ServerLoginResult byte

// Bot ...
type Bot struct {
	gorm.Model

	//BotID int `gorm:"primaryKey, autoIncrement:true"`
	//
	// UserID ower bot
	//
	UserID int
	//
	// UID уникальный индификатор бота (конконтинация SocialName + LoginID)
	//
	UID string `gorm:"unique"`
	//
	//InnerID игровой ID (используется внутри игры - публичный)
	//
	InnerID int32
	//
	// Result - ответ авторизации сервера (пока тут как строка)
	//
	Result string
	//
	//Ник
	//
	Name string
	//
	//Фото
	//
	Avatar string
	//
	//Ссылка на соц.сеть
	//
	Profile string
	//
	//Баланс
	//
	Balance int32
	//
	// IsBonus - флаг указывающий есть ли доступный бонус для получения
	//
	IsBonus bool
	//
	// BonusDay - счетчик дней бонуса [0..7]
	//
	BonusDay int32
	//
	// IsError -
	//
	IsError bool

	//
	// LastUseDay
	//
	LastUseDay string

	//
	// Logger - массив записей
	// `gorm:"type:text"`
	// `gorm:"foreignKey:LineID"`
	// PRELOAD:true
	// `gorm:"foreignKey:line_id, auto_preload:true"`
	//`gorm:"auto_preload:true"`
	Logger []LoggerLine

	//
	// Login - структура для авторизации на сервере. SocialCode - тип фрейма.
	// Если SocialCode === (255) - ОШИБКА ПАРСИНГА ФРЕЙМА
	//
	parser.LoginParams
}

// NewBot конструктор для Bot
func NewBot(url string) *Bot {

	bot := &Bot{

		InnerID: 0,
		Result:  "WAIT_AUTHORIZATION",
		Balance: 0,
		Name:    "No name",
		Avatar:  "ava.jpg",
		Profile: "",

		IsBonus:  false,
		BonusDay: 0,

		LoginParams: *parser.NewLoginParams(url),
		Logger:      make([]LoggerLine, 0),
		LastUseDay:  time.Now().Format(until.TIME_FORMAT),
	}

	if bot.SocialCode == 255 {
		bot.IsError = true
		bot.LogERROR("NewLoginParams", "Не удалось определить тип фрейма")
		return bot
	}

	bot.UID = bot.GetUID()
	return bot
}

// NewBotWhitProxy ...
func NewBotWhitProxy(url string, proxy string) *Bot {

	bot := &Bot{

		InnerID: 0,
		Result:  "WAIT_AUTHORIZATION",
		Balance: 0,
		Name:    "No name",
		Avatar:  "ava.jpg",
		Profile: "",

		IsBonus:  false,
		BonusDay: 0,

		LoginParams: *parser.NewLoginParams(url),

		Logger: make([]LoggerLine, 0),
	}

	if bot.SocialCode == 255 {
		bot.IsError = true
		bot.LogERROR("NewLoginParams", "Не удалось определить тип фрейма")
		return bot
	}

	bot.UID = bot.GetUID()

	return bot
}

// ToString ...
func (bot *Bot) ToString() string {
	return fmt.Sprintf("Bot:\n Name: %v\n InnerID:  %v\n Balance: %v\n Result: %v\n IsBonus: %v\n Bonus: %v\n Avatar: %v\n Profile: %v\n LoginParams:\n Social: %v\n LoginID: %v\n Token: %v\n Token2:%v\n",
		bot.Name,
		bot.InnerID,
		bot.Balance,
		bot.Result,
		bot.IsBonus,
		bot.BonusDay,
		bot.Avatar,
		bot.Profile,
		bot.SocialName,
		bot.LoginID,
		bot.Token,
		bot.Token2,
	)
}
func (bot *Bot) PrintLog() {
	fmt.Println("Logger: ")
	for i, v := range bot.Logger {
		fmt.Printf("\t%d) %v\n", i, v)
	}
}

type LogType int8

const (
	INFO  = 0
	ERROR = 1
)

func (bot *Bot) LogINFO(methodName, msg string)  { bot.Log(INFO, methodName, msg) }
func (bot *Bot) LogERROR(methodName, msg string) { bot.Log(ERROR, methodName, msg) }
func (bot *Bot) Log(t LogType, methodName, msg string) {

	switch t {
	case INFO:
		msg = fmt.Sprintf("INFO >> (%s) %s", methodName, msg)
	case ERROR:
		msg = fmt.Sprintf("ERROR >> (%s) %s", methodName, msg)
	}

	fmt.Println(msg)
	bot.Logger = append(bot.Logger, LoggerLine{
		Line: msg,
	})
}
