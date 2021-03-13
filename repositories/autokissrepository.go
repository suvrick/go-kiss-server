package repositories

import (
	"github.com/suvrick/go-kiss-server/model"
	"gorm.io/gorm"
)

type AutoKissRepository struct {
	db *gorm.DB
}

func NewAutoKissRepository(db *gorm.DB) *AutoKissRepository {
	repo := &AutoKissRepository{
		db: db,
	}

	return repo
}

// Add ...
func (repo *AutoKissRepository) Add(userID int) error {
	user := &model.KissUser{
		UserID:    userID,
		IsTrial:   true,
		KissCount: 0,
	}

	return repo.db.Create(user).Error
}

//
// Save ...
func (repo *AutoKissRepository) Save(user *model.KissUser) error {
	return repo.db.Save(user).Error
}

// Update ...
func (repo *AutoKissRepository) Update(user *model.KissUser) error {
	return repo.db.Table("kiss_users").Where("user_id = ?", user.UserID).Save(user).Error
}

// FindByID ...
func (repo *AutoKissRepository) FindByID(userID int) (*model.KissUser, error) {
	user := &model.KissUser{}
	err := repo.db.Table("kiss_users").Where("user_id = ?", userID).First(user).Error
	return user, err
}

// All ...
func (repo *AutoKissRepository) All() ([]model.KissUser, error) {
	users := make([]model.KissUser, 0)
	result := repo.db.Table("kiss_users").Scan(&users)
	return users, result.Error
}
