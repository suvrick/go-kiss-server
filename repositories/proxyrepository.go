package repositories

import (
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/model"
	"gorm.io/gorm"
)

// ProxyRepository ...
type ProxyRepository struct {
	db *gorm.DB
}

// NewProxyRepository ...
func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	return &ProxyRepository{
		db: db,
	}
}

// Create ...
func (r *ProxyRepository) Create(p model.Proxy) (int, error) {
	result := r.db.Create(&p)
	return p.ID, result.Error
}

// Update ...
func (r *ProxyRepository) Update(p model.Proxy) error {
	return r.db.Save(&p).Error
}

// Free ...
func (r *ProxyRepository) Free() (model.Proxy, error) {
	p := model.Proxy{}
	result := r.db.Where("is_bad = ? AND is_busy = ?", false, false).First(&p)

	if result.Error == gorm.ErrRecordNotFound {
		return p, errors.ErrProxyListEmpty
	}

	p.IsBusy = true
	r.Update(p)
	return p, result.Error
}

// All ...
func (r *ProxyRepository) All() ([]model.Proxy, error) {
	proxies := make([]model.Proxy, 0)
	result := r.db.Table("proxies").Scan(&proxies)
	return proxies, result.Error
}

// DeleteAll ...
func (r *ProxyRepository) DeleteAll() error {
	return r.db.Table("proxies").Where("is_busy = ?", true).Delete(&model.Proxy{}).Error
}
