package repositories

import (
	"time"

	"github.com/suvrick/go-kiss-server/model"
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
func (repo *BotRepository) Add(bot model.Bot) (int, error) {
	result := repo.db.Create(&bot)
	return bot.ID, result.Error
}

// Update ...
func (repo *BotRepository) Update(bot model.Bot) error {
	bot.DateUse = time.Now().Format("2006-01-02")
	return repo.db.Save(&bot).Error
}

// All ...
func (repo *BotRepository) All(userID string) ([]*model.Bot, error) {
	bots := make([]*model.Bot, 0)
	result := repo.db.Table("bots").Where("user_id = ?", userID)
	result.Scan(&bots)
	return bots, result.Error
}

// Find ...
func (repo *BotRepository) Find(botID int) (model.Bot, error) {
	bot := model.Bot{}
	err := repo.db.Table("bots").Where("id = ?", botID).First(&bot).Error
	return bot, err
}

// Delete ...
func (repo *BotRepository) Delete(bot model.Bot) error {
	return repo.db.Delete(&bot).Error
}
