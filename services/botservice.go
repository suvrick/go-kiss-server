package services

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/ws"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/until"
)

// BotService ...
type BotService struct {
	userService   *UserService
	botRepository *repositories.BotRepository

	locker *sync.Mutex
}

// NewBotService ...
func NewBotService(repo *repositories.BotRepository, us *UserService) *BotService {
	return &BotService{
		userService:   us,
		botRepository: repo,

		locker: &sync.Mutex{},
	}
}

// Add ...
func (s *BotService) Add(c *gin.Context, url string) (*models.Bot, error) {

	s.locker.Lock()
	defer s.locker.Unlock()

	user := session.GetUser(c)

	if err := s.checkActualUser(user); err != nil {
		return nil, err
	}

	bot := models.NewBot(url)

	bot.UserID = user.ID

	if bot.UID == "" {
		return nil, errors.ErrParseBotInvalid
	}

	gs := ws.NewSocket(bot)
	gs.Go()

	_, err := s.botRepository.Add(*bot)

	if err != nil {
		return nil, err
	}

	user.BotsCount++
	if err := s.userService.UpdateUser(user); err != nil {
		return bot, err
	}

	return bot, nil
}

func (s *BotService) UpdateByID(c *gin.Context, id string) (*models.Bot, error) {
	user := session.GetUser(c)
	bot, err := s.botRepository.Find(id, user.ID)

	if err != nil {
		return nil, err
	}

	gs := ws.NewSocket(&bot)
	gs.Go()

	err = s.botRepository.Update(bot)

	return &bot, err
}

// UpdateBot ...
func (s *BotService) UpdateBot(bot models.Bot) {
	s.botRepository.Update(bot)
}

// All ...
func (s *BotService) All(c *gin.Context) ([]*models.Bot, error) {
	user := session.GetUser(c)
	return s.botRepository.All(user.ID)
}

// AllByUserID bots by userID ...
func (s *BotService) AllByUserID(userID int) ([]*models.Bot, error) {
	return s.botRepository.All(userID)
}

// Delete ...
func (s *BotService) Delete(c *gin.Context) error {

	s.locker.Lock()
	defer s.locker.Unlock()

	user := session.GetUser(c)
	if user.ID == 0 {
		return errors.ErrNotAuthenticated
	}

	botUID := c.Param("botID")

	bot, err := s.botRepository.Find(botUID, user.ID)

	if bot.UserID != user.ID || err != nil {
		return errors.ErrRecordNotFound
	}

	if err := s.botRepository.Delete(bot); err != nil {
		return errors.ErrRecordNotFound
	}

	user.BotsCount--
	s.userService.UpdateUser(user)

	return nil
}

// CheckActualUser ...
func (s *BotService) checkActualUser(user *model.User) error {

	userDate, _ := time.Parse(until.TIME_FORMAT, user.Date)
	nowDate := time.Now()

	if userDate.Unix() < nowDate.Unix() {
		return errors.ErrNeedTime
	}

	if user.BotsCount >= user.Limit {
		return errors.ErrLimitBot
	}

	return nil
}
