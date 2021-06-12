package repositories

import (
	"github.com/suvrick/go-kiss-server/model"
	"gorm.io/gorm"
)

// UserRepository ...
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ...
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create ...
func (r *UserRepository) Create(u model.User) (int, error) {
	result := r.db.Create(&u)
	return u.ID, result.Error
}

// UpdateUser ...
func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// FindByToken ...
func (r *UserRepository) FindByToken(key string) (*model.User, error) {
	user := &model.User{}
	result := r.db.Table("users").Where("token = ?", key).First(user)
	return user, result.Error
}

// FindByID ...
func (r *UserRepository) FindByID(userID int) (*model.User, error) {
	user := &model.User{}
	result := r.db.Table("users").Where("id = ?", userID).First(user)
	return user, result.Error
}

// FindByEmailAndPass ...
func (r *UserRepository) FindByEmailAndPass(email, password string) (*model.User, error) {
	user := &model.User{}
	result := r.db.Table("users").Where("email = ? AND password = ?", email, password).First(user)
	return user, result.Error
}

// All ...
func (r *UserRepository) All() ([]model.User, error) {
	users := make([]model.User, 0)
	result := r.db.Table("users").Where("ID != ?", 0).Scan(&users)
	return users, result.Error
}
