package service

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/repository/postgres"
)

type locationService struct {
	repo postgres.LocationRepository
}

func newLocationService(r postgres.LocationRepository) *locationService {
	return &locationService{
		repo: r,
	}
}

func (s *locationService) GetAll() (model.Locations, error) {
	locations, err := s.repo.FindAll()
	if err != nil {
		return locations, err
	}
	return locations, err
}

func (s *locationService) GetById(id string) (model.Location, error) {
	location, err := s.repo.FindById(id)
	if err != nil {
		return location, err
	}
	return location, nil
}

func (s *locationService) GetRating(id string) (float32, error) {
	rating, err := s.repo.FindRating(id)
	if err != nil {
		return rating, err
	}
	return rating, err
}

func (s *locationService) Create(loc model.Location) (model.Location, error) {
	location, err := s.repo.Insert(loc)
	if err != nil {
		return location, err
	}
	return location, err
}
