package models

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-server/game/parser"
)

type ServerLoginResult byte

// Bot ...
type Bot struct {
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
	//
	Logger []string `gorm:"-"`

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
		Logger:      make([]string, 0),
		LastUseDay:  time.Now().Format("2006-01-02"),
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

		Logger: make([]string, 0),
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
	return fmt.Sprintf("****Bot****\n Name: %v\n InnerID:  %v\n Balance: %v\n Result: %v\n IsBonus: %v\n Bonus: %v\n Avatar: %v\n Profile: %v\n LoginParams:\n Social: %v\n LoginID: %v\n Token: %v\n Token2:%v\n",
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
		bot.Logger = append(bot.Logger, fmt.Sprintf("INFO (%s) %s", methodName, msg))
	case ERROR:
		bot.Logger = append(bot.Logger, fmt.Sprintf("ERROR >> (%s) %s", methodName, msg))
	}
}
