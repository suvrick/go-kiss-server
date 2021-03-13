package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-core/models"
	"github.com/suvrick/go-kiss-core/parser"
	"github.com/suvrick/go-kiss-core/ws"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/until"
)

// BotService ...
type BotService struct {
	botRepository *repositories.BotRepository
	userService   *UserService
	proxyService  *ProxyService
}

// NewBotService ...
func NewBotService(repo *repositories.BotRepository, us *UserService, ps *ProxyService) *BotService {
	return &BotService{
		botRepository: repo,
		userService:   us,
		proxyService:  ps,
	}
}

// Add ...
func (s *BotService) Add(c *gin.Context, url string) (*model.Bot, error) {

	userID, user, err := until.GetUserFromContext(c)
	if err != nil {
		return nil, err
	}

	if err := s.CheckActualUser(userID, user); err != nil {
		return nil, err
	}

	loginData := parser.NewLoginParams(url)
	bot := model.Bot{
		LoginParams: loginData,
	}

	bot.UserID = userID

	botID, err := s.botRepository.Add(bot)
	if err != nil {
		return nil, err
	}

	bot.ID = botID
	bot.DateUse = time.Now().Format("2006-01-02")

	user.BotsCount++
	s.userService.UpdateUser(user)

	if _, err := s.InGame(&bot); err != nil {
		return nil, err
	}

	return &bot, nil
}

// UpdateBot ...
func (s *BotService) UpdateBot(bot model.Bot) {

	s.botRepository.Update(bot)
}

// CheckActualUser ...
func (s *BotService) CheckActualUser(userID string, user model.User) error {

	userDate, _ := time.Parse("2006-01-02", user.Date)
	nowDate := time.Now()

	if userDate.Unix() < nowDate.Unix() {
		return errors.ErrNeedTime
	}

	if user.BotsCount >= user.Limit {
		return errors.ErrLimitBot
	}

	return nil
}

// InGame ....
func (s *BotService) InGame(bot *model.Bot) (*model.Bot, error) {

	p, err := s.proxyService.Free()

	if err != nil {
		return bot, err
	}

	botGame := models.NewBotWhitProxy(bot.LoginParams.FrameURL, p.URL)
	g := ws.NewSocket(botGame)
	g.Go()

	fmt.Println("Get post update bot:", botGame)
	bot.Result = botGame.Result

	if botGame.IsError {
		fmt.Println(botGame.Error)
		return bot, err
	}

	s.UpdateBot(*bot)
	return bot, nil
}

// All ...
func (s *BotService) All(c *gin.Context) ([]*model.Bot, error) {

	userID, user, err := until.GetUserFromContext(c)
	if err != nil {
		return nil, err
	}

	userDate, _ := time.Parse("2006-01-02", user.Date)
	nowDate := time.Now()

	if userDate.Unix() < nowDate.Unix() {
		return nil, errors.ErrNeedTime
	}

	return s.botRepository.All(userID)
}

// Delete ...
func (s *BotService) Delete(c *gin.Context) error {

	userID, _, err := until.GetUserFromContext(c)
	if err != nil {
		return err
	}

	botIDStr := c.Param("botID")

	botID, err := strconv.Atoi(botIDStr)
	if err != nil {
		return errors.ErrInvalidParam
	}

	bot, err := s.Find(botID)

	if bot.UserID != userID {
		return errors.ErrInvalidParam
	}

	return s.botRepository.Delete(bot)
}

// Find ...
func (s *BotService) Find(botID int) (model.Bot, error) {
	return s.botRepository.Find(botID)
}
