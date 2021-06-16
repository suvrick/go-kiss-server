package repositories

import (
	"sync"
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/until"
	"gorm.io/gorm"
)

// BotRepository ...
type BotRepository struct {
	db *gorm.DB

	locker *sync.Mutex
}

// NewBotRepository ...
func NewBotRepository(db *gorm.DB) *BotRepository {
	return &BotRepository{
		db: db,

		locker: &sync.Mutex{},
	}
}

// Add ...
func (repo *BotRepository) Add(bot *models.Bot) (string, error) {

	repo.locker.Lock()
	defer repo.locker.Unlock()

	result := repo.db.Table("bots").Create(bot)
	return bot.UID, result.Error
}

// Update ...
func (repo *BotRepository) Update(bot *models.Bot) error {

	repo.locker.Lock()
	defer repo.locker.Unlock()

	nb := &models.Bot{}
	repo.db.First(nb, "uid = ?", bot.UID)

	repo.db.Unscoped().Table("logger_lines").Where("bot_id", nb.ID).Delete(&[]models.LoggerLine{})

	bot.LastUseDay = time.Now().Format(until.TIME_FORMAT)
	return repo.db.Preload("Logger").Save(bot).Error
}

// All ...
func (repo *BotRepository) All(userID int) ([]*models.Bot, error) {

	repo.locker.Lock()
	defer repo.locker.Unlock()

	bots := make([]*models.Bot, 0)
	result := repo.db.Preload("Logger").Find(&bots, "user_id = ?", userID)
	return bots, result.Error
}

// Find ...
func (repo *BotRepository) Find(botUID string, userID int) (*models.Bot, error) {

	repo.locker.Lock()
	defer repo.locker.Unlock()

	bot := &models.Bot{}
	err := repo.db.Preload("Logger").First(bot, "uid = ? AND user_id = ?", botUID, userID).Error
	return bot, err
}

// Delete ...
func (repo *BotRepository) Delete(bot *models.Bot) error {

	repo.locker.Lock()
	defer repo.locker.Unlock()

	repo.db.Unscoped().Table("logger_lines").Where("bot_id", bot.ID).Delete(&[]models.LoggerLine{})
	return repo.db.Unscoped().Where("uid = ?", bot.UID).Delete(&models.Bot{}).Error
}
