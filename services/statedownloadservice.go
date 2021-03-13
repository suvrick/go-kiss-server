package services

import (
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
)

type StateDownloadService struct {
	repo *repositories.StateDownloadRepository
}

func NewStateDownloadService(r *repositories.StateDownloadRepository) *StateDownloadService {
	s := &StateDownloadService{
		repo: r,
	}

	return s
}

func (s *StateDownloadService) AddDownloadState(ip string) error {
	return s.repo.Add(ip)
}

func (s *StateDownloadService) FindDownloadState(ip string) (*model.KissState, error) {
	return s.repo.Find(ip)
}

func (s *StateDownloadService) UpdateDownloadState(state *model.KissState) error {
	return s.repo.Update(state)
}

func (s *StateDownloadService) AllDownloadState() ([]model.KissState, error) {
	return s.repo.All()
}
