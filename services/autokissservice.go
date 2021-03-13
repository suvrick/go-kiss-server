package services

import (
	"encoding/binary"

	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
)

type AutoKissService struct {
	repo *repositories.AutoKissRepository
}

const (
	TRIAL_KISS_COUNT = 30
)

func NewAutoKissService(r *repositories.AutoKissRepository) *AutoKissService {
	service := &AutoKissService{
		repo: r,
	}

	return service
}

func (s *AutoKissService) Do(ID int, t uint16, data []byte) *model.KissResponse {

	res := &model.KissResponse{}

	switch t {
	case 29:

		if len(data) < 14 {
			return res
		}

		leaderID := int(binary.LittleEndian.Uint32(data[6:10]))
		rolledID := int(binary.LittleEndian.Uint32(data[10:14]))

		if leaderID != ID && rolledID != ID {
			return res
		}

		//log.Printf("leaderID: %v rolledID: %v", leaderID, rolledID)

		u1, _ := s.repo.FindByID(ID)

		if u1.UserID != 0 {

			if u1.KissCount > TRIAL_KISS_COUNT && u1.IsTrial {
				res.Status = 403
				return res
			}

			res.Status = 200
			res.Code = 29
			res.Data = []interface{}{1}
			res.Delay = 7000

			u1.KissCount++
			s.repo.Update(u1)

		}

		return res

	case 28:

		if len(data) < 10 {
			return res
		}

		leaderID := int(binary.LittleEndian.Uint32(data[6:10]))

		if leaderID != ID {
			return res
		}

		user, _ := s.repo.FindByID(ID)

		if user.UserID != 0 {

			if user.KissCount > TRIAL_KISS_COUNT && user.IsTrial {
				res.Status = 403
				return res
			}

			res.Code = 28
			res.Data = []interface{}{0}
			res.Delay = 5000
		}

		return res

	case 308:

		if len(data) < 14 {
			return res
		}

		kickID := int(binary.LittleEndian.Uint32(data[6:10]))
		if kickID != ID {
			return res
		}

		user, _ := s.repo.FindByID(kickID)

		if user.UserID != 0 {

			if user.KissCount > TRIAL_KISS_COUNT && user.IsTrial {
				res.Status = 403
				return res
			}

			res.Code = 30
			res.Data = []interface{}{kickID}
			res.Delay = 4000
		}

		return res
	}

	return res
}

func (s *AutoKissService) AddKissUser(userID int) error {
	return s.repo.Add(userID)
}

func (s *AutoKissService) FindByIDKissUser(userID int) (*model.KissUser, error) {
	return s.repo.FindByID(userID)
}

func (s *AutoKissService) SaveKissUser(user *model.KissUser) error {
	return s.repo.Save(user)
}

func (s *AutoKissService) UpdateKissUser(user *model.KissUser) error {
	return s.repo.Update(user)
}

func (s *AutoKissService) AllKissUser() ([]model.KissUser, error) {
	return s.repo.All()
}
