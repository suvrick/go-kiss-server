package repositories

import (
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"gorm.io/gorm"
)

// BotRepository ...
type BotRepository struct {
	db *gorm.DB
}

// NewBotRepository ...
func NewBotRepository(db *gorm.DB) *BotRepository {
	return &BotRepository{
		db: db,
	}
}

// Add ...
func (repo *BotRepository) Add(bot models.Bot) (string, error) {
	result := repo.db.Table("bots").Create(&bot)
	return bot.UID, result.Error
}

// Update ...
func (repo *BotRepository) Update(bot models.Bot) error {
	bot.LastUseDay = time.Now().Format("2006-01-02")
	return repo.db.Save(&bot).Error
}

// All ...
func (repo *BotRepository) All(userID int) ([]*models.Bot, error) {
	bots := make([]*models.Bot, 0)
	result := repo.db.Table("bots").Where("user_id = ?", userID).Scan(&bots)
	return bots, result.Error
}

// Find ...
func (repo *BotRepository) Find(botUID string, userID int) (models.Bot, error) {
	bot := models.Bot{}
	err := repo.db.Table("bots").Where("uid = ? AND user_id = ?", botUID, userID).First(&bot).Error
	return bot, err
}

// Delete ...
func (repo *BotRepository) Delete(bot models.Bot) error {
	return repo.db.Delete(&bot).Error
}
