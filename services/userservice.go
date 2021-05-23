package services

import (
	"sync"
	"time"

	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/until"
)

type UserService struct {
	userRepo *repositories.UserRepository
	locker   *sync.Mutex
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	s := &UserService{
		userRepo: repo,
		locker:   &sync.Mutex{},
	}

	return s
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	return s.userRepo.FindByEmailAndPass(email, password)
}

func (s *UserService) Register(email, password string) (int, error) {

	//default role player
	//after deploy call cmd
	//update users set role = 'admin' where id = 1;

	u := model.User{
		Email:    email,
		Password: password,
		Role:     "player",
		Date:     time.Now().Format(until.TIME_FORMAT),
	}

	return s.userRepo.Create(u)
}

func (s *UserService) FindUserByID(userID int) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateUser(user *model.User) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	session.Accounts[user.Token] = user
	return s.userRepo.UpdateUser(user)
}

func (s *UserService) AllUser() ([]model.User, error) {
	return s.userRepo.All()
}
