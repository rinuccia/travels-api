package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
)

type visitRepo struct {
	*sqlx.DB
}

func newVisitRepo(db *sqlx.DB) *visitRepo {
	return &visitRepo{db}
}

func (r *visitRepo) FindAll(id string) (model.UserVisits, error) {
	var userId uint32
	visit := model.UserVisit{}
	visits := model.UserVisits{}
	row := r.QueryRow("SELECT user_id FROM users WHERE user_id = $1", id)
	err := row.Scan(&userId)
	if err != nil {
		return visits, apperrors.ErrRecordNotFound
	}
	query := `
			SELECT locations.place, locations.country, visits.visited_at, visits.mark
			FROM users 
				JOIN visits  
					ON users.user_id = visits.user_id 
				JOIN locations 
					ON locations.location_id = visits.location_id 
			WHERE users.user_id = $1
			ORDER BY visits.visited_at`
	rows, err := r.Query(query, id)
	if err != nil {
		return visits, err
	}
	for rows.Next() {
		err = rows.Scan(&visit.Place, &visit.Country, &visit.VisitedAt, &visit.Mark)
		if err != nil {
			return visits, err
		}
		visits.Visits = append(visits.Visits, visit)
	}
	return visits, err
}

func (r *visitRepo) Insert(visit model.Visit) (model.Visit, error) {
	query := "INSERT INTO visits (visit_id, location_id, user_id, visited_at, mark) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.Exec(query, visit.VisitId, visit.LocationId, visit.UserId, visit.VisitedAt, visit.Mark)
	if err != nil {
		return visit, apperrors.ErrIncorrectQuery
	}
	return visit, err
}

func (r *visitRepo) DeleteById(id string) error {
	res, err := r.Exec("DELETE FROM visits WHERE visit_id = $1", id)
	if err != nil {
		return err
	}
	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		return apperrors.ErrRecordNotFound
	}
	return err
}
