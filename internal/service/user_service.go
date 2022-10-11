package service

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/repository/postgres"
)

type userService struct {
	repo postgres.UserRepository
}

func newUserService(r postgres.UserRepository) *userService {
	return &userService{
		repo: r,
	}
}

func (s *userService) GetById(id string) (model.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *userService) Create(user model.User) (model.User, error) {
	var err error
	user, err = s.repo.Insert(user)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *userService) Update(id string, user model.User) error {
	err := s.repo.Update(id, user)
	if err != nil {
		return err
	}
	return err
}
