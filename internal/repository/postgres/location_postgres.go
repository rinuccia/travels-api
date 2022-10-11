package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
)

type locationRepo struct {
	*sqlx.DB
}

func newLocationRepo(db *sqlx.DB) *locationRepo {
	return &locationRepo{db}
}

func (r *locationRepo) FindAll() (model.Locations, error) {
	query := "SELECT * FROM locations"
	location := model.Location{}
	locations := model.Locations{}
	rows, err := r.Query(query)
	if err != nil {
		return locations, err
	}
	for rows.Next() {
		err = rows.Scan(&location.LocationId, &location.Place, &location.Country)
		if err != nil {
			return locations, err
		}
		locations.List = append(locations.List, location)
	}
	return locations, err
}

func (r *locationRepo) FindById(id string) (model.Location, error) {
	query := "SELECT * FROM locations WHERE location_id = $1"
	location := model.Location{}
	row := r.QueryRow(query, id)
	err := row.Scan(&location.LocationId, &location.Place, &location.Country)
	if err != nil {
		return location, apperrors.ErrRecordNotFound
	}
	return location, err
}

func (r *locationRepo) FindRating(id string) (float32, error) {
	var locationId int
	row := r.QueryRow("SELECT location_id FROM locations WHERE location_id = $1", id)
	err := row.Scan(&locationId)
	if err != nil {
		return 0, apperrors.ErrRecordNotFound
	}
	query := `
			SELECT ROUND(sum * 1.0 / count, 2) AS avg 
			FROM (SELECT locations.location_id, SUM(visits.mark) AS sum, COUNT(visits.mark) AS count 
				  FROM locations 
					  LEFT JOIN visits 
						  ON locations.location_id = visits.location_id 
				  GROUP BY locations.location_id 
				  HAVING locations.location_id = $1) AS abc`
	var rating float32
	row = r.QueryRow(query, id)
	_ = row.Scan(&rating)
	return rating, err
}

func (r *locationRepo) Insert(location model.Location) (model.Location, error) {
	query := "INSERT INTO locations (location_id, place, country) VALUES ($1, $2, $3)"
	_, err := r.Exec(query, location.LocationId, location.Place, location.Country)
	if err != nil {
		return location, apperrors.ErrIncorrectQuery
	}

	return location, err
}
