package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/rinuccia/travels-api/internal/model"
)

type (
	UserRepository interface {
		// FindById user in DB.
		FindById(id string) (model.User, error)

		// Insert user with given credentials in DB.
		Insert(u model.User) (model.User, error)

		// Update user in DB.
		Update(id string, u model.User) error
	}

	LocationRepository interface {
		// FindAll locations in DB.
		FindAll() (model.Locations, error)

		// FindById user in DB.
		FindById(id string) (model.Location, error)

		// FindRating location by id.
		FindRating(id string) (float32, error)

		// Insert location with given credentials in DB.
		Insert(location model.Location) (model.Location, error)
	}

	VisitRepository interface {
		// FindAll user visits by id in DB.
		FindAll(id string) (model.UserVisits, error)

		// Insert new visit in DB.
		Insert(visit model.Visit) (model.Visit, error)

		// DeleteById user visit in DB.
		DeleteById(id string) error
	}
)

type Repository struct {
	UserRepository
	LocationRepository
	VisitRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		newUserRepo(db),
		newLocationRepo(db),
		newVisitRepo(db),
	}
}
