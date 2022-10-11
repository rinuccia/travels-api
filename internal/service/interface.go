package service

import (
	"github.com/rinuccia/travels-api/internal/model"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type (
	User interface {
		// GetById user.
		GetById(id string) (model.User, error)

		// Create new user.
		Create(user model.User) (model.User, error)

		// Update user by id.
		Update(id string, user model.User) error
	}

	Location interface {
		// GetAll locations.
		GetAll() (model.Locations, error)

		// GetById location.
		GetById(id string) (model.Location, error)

		// GetRating location by id.
		GetRating(id string) (float32, error)

		// Create new location.
		Create(loc model.Location) (model.Location, error)
	}

	Visit interface {
		// GetAll user visits by id.
		GetAll(id string) (model.UserVisits, error)

		// Create new visit.
		Create(visit model.Visit) (model.Visit, error)

		// DeleteById user visit.
		DeleteById(id string) error
	}
)
