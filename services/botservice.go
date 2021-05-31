package services

import (
	"sync"
	"time"

	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/packets/encode"
	"github.com/suvrick/go-kiss-server/game/ws"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
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
func (s *BotService) Add(user *model.User, url string) (*models.Bot, error) {

	s.locker.Lock()
	defer s.locker.Unlock()

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

	_, err := s.botRepository.Add(bot)

	if err != nil {
		return nil, err
	}

	user.BotsCount++
	if err := s.userService.UpdateUser(user); err != nil {
		return bot, err
	}

	return bot, nil
}

func (s *BotService) UpdateByID(botUID string, user *model.User) (*models.Bot, error) {

	bot, err := s.botRepository.Find(botUID, user.ID)

	if err != nil {
		return nil, err
	}

	gs := ws.NewSocket(bot)
	gs.Go()

	err = s.botRepository.Update(bot)

	return bot, err
}

func (s *BotService) SendPrize(botUID string, user *model.User, prize encode.ClientPacket) (*models.Bot, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	bot, err := s.botRepository.Find(botUID, user.ID)

	if bot.UserID != user.ID || err != nil {
		return nil, errors.ErrRecordNotFound
	}

	gs := ws.NewSocketWithPrize(bot, &prize)
	gs.Go()

	err = s.botRepository.Update(bot)

	return bot, err
}

func (s *BotService) SendPrize2(botUID string, user *model.User, prize encode.ClientPacket) (*models.Bot, error) {

	bot, err := s.botRepository.Find(botUID, user.ID)

	if bot.UserID != user.ID || err != nil {
		return nil, errors.ErrRecordNotFound
	}

	gs := ws.NewSocketWithPrize(bot, &prize)
	gs.Go()

	err = s.botRepository.Update(bot)

	return bot, err
}

// UpdateBot ...
func (s *BotService) UpdateBot(bot *models.Bot) error {
	return s.botRepository.Update(bot)
}

// All ...
func (s *BotService) All(user *model.User) ([]*models.Bot, error) {
	return s.botRepository.All(user.ID)
}

// AllByUserID bots by userID ...
func (s *BotService) AllByUserID(userID int) ([]*models.Bot, error) {
	return s.botRepository.All(userID)
}

// Delete ...
func (s *BotService) Delete(botUID string, user *model.User) error {

	s.locker.Lock()
	defer s.locker.Unlock()

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
