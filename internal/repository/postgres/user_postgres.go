package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
)

type userRepo struct {
	*sqlx.DB
}

func newUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) FindById(id string) (model.User, error) {
	query := "SELECT * FROM users WHERE user_id = $1"
	user := model.User{}
	row := r.QueryRow(query, id)
	err := row.Scan(&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Gender)
	if err != nil {
		return user, apperrors.ErrRecordNotFound
	}

	return user, err
}

func (r *userRepo) Insert(user model.User) (model.User, error) {
	query := "INSERT INTO users (user_id, email, first_name, last_name, gender) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.Exec(query, user.UserId, user.Email, user.FirstName, user.LastName, user.Gender)
	if err != nil {
		return user, apperrors.ErrIncorrectQuery
	}

	return user, err
}

func (r *userRepo) Update(id string, u model.User) error {
	query := "UPDATE users SET email = $1, first_name = $2, last_name = $3, gender = $4 WHERE user_id = $5"

	res, err := r.Exec(query, u.Email, u.FirstName, u.LastName, u.Gender, id)
	if err != nil {
		return apperrors.ErrIncorrectQuery
	}
	if rowsAff, err := res.RowsAffected(); rowsAff == 0 && err == nil {
		return apperrors.ErrRecordNotFound
	}
	return err
}
