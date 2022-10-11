package service

import "github.com/rinuccia/travels-api/internal/repository/postgres"

type Service struct {
	User
	Location
	Visit
}

func NewService(repos *postgres.Repository) *Service {
	return &Service{
		newUserService(repos.UserRepository),
		newLocationService(repos.LocationRepository),
		newVisitService(repos.VisitRepository),
	}
}
