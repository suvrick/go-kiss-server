// Package parser ...
// Пакет для парсинга фреймов в структуру LoginParams
package parser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// ParseParamKey ..
type ParseParamKey struct {
	SType     string `json:"stype"`
	STypeCode int16  `json:"stypecode"`
	ID        string `json:"id"`
	Token     string `json:"token"`
	Token2    string `json:"token2"`
}

// LoginParams struct
// Отдаем наружу
type LoginParams struct {
	SocialCode int16
	//двухбуквеный код социалки
	SocialName string
	LoginID    int64
	Token      string
	Token2     string
	FrameURL   string `json:"-"`
}

const (
	debug = false
)

var parseParams []ParseParamKey

//go:embed config.json
var pathParserConfig []byte

// init ...
func init() {
	Initialize()
}

// Initialize ...
// Загрузка ключей для парсинга строки
func Initialize() error {

	parseParams = []ParseParamKey{}
	if err := json.Unmarshal(pathParserConfig, &parseParams); err != nil {
		log.Fatalf("Parser [init] >> %s", err.Error())
		return err
	}

	return nil
}

// NewLoginParams ...
//
// Парсим url в LoginParams структуру
//
// LoginParams.SocialCode - тип регистрации (int16)
//
// Если неудача -> LoginParams.SocialCode == 255
//
// LoginParams.ErrorMsg - сообщения об ошибки
//
func NewLoginParams(input string) *LoginParams {

	params := &LoginParams{
		FrameURL:   input,
		SocialName: "nn",
		SocialCode: 255,
	}

	// Удаляем пробелы и спец.символы
	input = strings.TrimSpace(input)
	input = strings.Replace(input, "\r", "", -1)

	if len(input) == 0 {
		return params
	}

	// Пытаемся получить map[ключ]=значения, query элементов URL
	query, err := url.ParseQuery(input)
	if err != nil {
		return params
	}

	/*
		Тут уже начинается хардкор :D
		Пытаемся как-нибудь определить тип социалки
	*/

	for _, p := range parseParams {
		if strings.Contains(input, p.ID) &&
			strings.Contains(input, p.Token) &&
			strings.Contains(input, p.Token2) {

			p.Debug(input)

			strID := query.Get(p.ID)
			loginID, err := strconv.ParseInt(strID, 10, 64)
			if err != nil || loginID == 0 {
				return params
			}

			params.LoginID = loginID
			params.SocialName = p.SType
			params.SocialCode = p.STypeCode
			params.Token = query.Get(p.Token)
			params.Token2 = query.Get(p.Token2)

			return params
		}
	}

	// Возращаем дефолт 0xFF
	return params
}

// GetUID возращаем уникальный индификатор для Player
func (lp *LoginParams) GetUID() string {
	return fmt.Sprintf("%s%d", lp.SocialName, lp.LoginID)
}

// ToString конвертирует структуру LoginParams в строку
func (lp *LoginParams) ToString() string {
	return fmt.Sprintf(" Social type: %s\n LoginID: %d\n Token: %s\n Token2: %s\n",
		lp.SocialName,
		lp.LoginID,
		lp.Token,
		lp.Token2,
	)
}

// Debug выводим в консоль ParseParamKey
//
// Проверяем вхождения ключей в строке
//
// Должен быть включен флаг debug
func (p ParseParamKey) Debug(input string) {
	if debug {
		log.Printf("TYPE %s", p.SType)
		log.Printf("TYPE CODE %d", p.STypeCode)
		log.Printf("ID %s, %v", p.ID, strings.Contains(input, p.ID))
		log.Printf("TOKEN %s, %v", p.Token, strings.Contains(input, p.Token))
		log.Printf("TOKEN2 %s, %v", p.Token2, strings.Contains(input, p.Token2))
		log.Printf("%s\n", input)
	}
}
