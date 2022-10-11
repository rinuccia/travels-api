package service

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/repository/postgres"
)

type visitService struct {
	repo postgres.VisitRepository
}

func newVisitService(r postgres.VisitRepository) *visitService {
	return &visitService{
		repo: r,
	}
}

func (s *visitService) GetAll(id string) (model.UserVisits, error) {
	visits, err := s.repo.FindAll(id)
	if err != nil {
		return visits, err
	}
	return visits, err
}

func (s *visitService) Create(visit model.Visit) (model.Visit, error) {
	v, err := s.repo.Insert(visit)
	if err != nil {
		return v, err
	}
	return v, err
}

func (s *visitService) DeleteById(id string) error {
	err := s.repo.DeleteById(id)
	return err
}
