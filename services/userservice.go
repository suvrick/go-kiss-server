package services

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/session"
)

// UserService ...
type UserService struct {
	userRepository *repositories.UserRepository
}

// NewUserService ...
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

// Create ...
func (s *UserService) Create(login, password string) (int, error) {

	u := model.User{
		Email:    login,
		Password: password,
		Date:     time.Now().Format("2006-01-02"),
	}

	return s.userRepository.Create(u)
}

// Login ...
func (s *UserService) Login(c *gin.Context, login, password string) (model.User, error) {

	u := model.User{
		Email:    login,
		Password: password,
		Date:     time.Now().Format("2006-01-02"),
	}

	user, err := s.userRepository.FindByEmailAndPass(login, password)
	if err != nil {
		return u, err
	}

	user.Token = s.GetMD5Hash(u.Email, u.Password)

	host := strings.Split(c.Request.Host, ":")[0]
	c.SetCookie("token", user.Token, 60*60*24, "/", host, false, false)
	session.Accounts[user.Token] = user

	return user, s.UpdateUser(user)

}

// FindByKey ...
func (s *UserService) FindByKey(key string) (model.User, error) {
	return s.userRepository.FindByKey(key)
}

// FindByID ...
func (s *UserService) FindByID(userID int) (model.User, error) {
	return s.userRepository.FindByID(userID)
}

// UpdateUser ...
func (s *UserService) UpdateUser(user model.User) error {
	session.Accounts[user.Token] = user
	return s.userRepository.UpdateUser(user)
}

// GetMD5Hash ...
func (s *UserService) GetMD5Hash(login, password string) string {
	now := time.Now().String()
	hash := md5.Sum([]byte(login + password + now))
	return hex.EncodeToString(hash[:])
}
