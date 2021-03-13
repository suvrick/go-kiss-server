package repositories

import (
	"time"

	"github.com/suvrick/go-kiss-server/model"
	"gorm.io/gorm"
)

type StateDownloadRepository struct {
	db *gorm.DB
}

func NewStateDowloadRepository(db *gorm.DB) *StateDownloadRepository {
	repo := &StateDownloadRepository{
		db: db,
	}

	return repo
}

// Add ...
func (repo *StateDownloadRepository) Add(ip string) error {
	state := &model.KissState{
		IP:   ip,
		Date: time.Now().Format(time.RFC822),
	}

	return repo.db.Create(state).Error
}

func (repo *StateDownloadRepository) Find(ip string) (*model.KissState, error) {
	state := &model.KissState{}
	err := repo.db.Table("kiss_states").Where("ip = ?", ip).First(state).Error
	return state, err
}

// Update ...
func (repo *StateDownloadRepository) Update(state *model.KissState) error {
	return repo.db.Table("kiss_states").Where("ip = ?", state.IP).Save(state).Error
}

// Add ...
func (repo *StateDownloadRepository) All() ([]model.KissState, error) {
	states := make([]model.KissState, 0)
	result := repo.db.Table("kiss_states").Scan(&states)
	return states, result.Error
}
