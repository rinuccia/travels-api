package postgres

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestLocationRepo_FindAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newLocationRepo(db)

	var testTable = []struct {
		name    string
		mock    func()
		want    model.Locations
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"location_id", "place", "country"}).
					AddRow(1, "Red Square", "RF").
					AddRow(2, "Eiffel Tower", "France").
					AddRow(3, "Grand Canyon", "USA")
				mock.ExpectQuery("SELECT (.+) FROM locations").WillReturnRows(rows)
			},
			want: model.Locations{
				List: []model.Location{
					{1, "Red Square", "RF"},
					{2, "Eiffel Tower", "France"},
					{3, "Grand Canyon", "USA"},
				},
			},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repository.FindAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLocationRepo_FindById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newLocationRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		id      string
		want    model.Location
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"location_id", "place", "country"}).
					AddRow("1", "Red Square", "RF")
				mock.ExpectQuery("SELECT (.+) FROM locations WHERE (.+)").WillReturnRows(rows)
			},
			id: "1",
			want: model.Location{
				LocationId: 1,
				Place:      "Red Square",
				Country:    "RF",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM locations WHERE (.+)").
					WithArgs("1")
			},
			id:      "1",
			wantErr: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repository.FindById(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLocationRepo_FindRating(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newLocationRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		id      string
		want    float32
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectQuery(`SELECT (.+) FROM locations`).WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"locations_id"}).AddRow(1))
				mock.ExpectQuery(`SELECT (.+) AS avg`).WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(4.5))
			},
			id:   "1",
			want: 4.5,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(`SELECT (.+) FROM locations`).WithArgs("1").
					WillReturnError(apperrors.ErrRecordNotFound)
			},
			id:      "1",
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repository.FindRating(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLocationRepo_Insert(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newLocationRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		input   model.Location
		want    model.Location
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO locations").
					WithArgs(1, "Machu Picchu", "Peru").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: model.Location{
				LocationId: 1,
				Place:      "Machu Picchu",
				Country:    "Peru",
			},
			want: model.Location{
				LocationId: 1,
				Place:      "Machu Picchu",
				Country:    "Peru",
			},
		},
		{
			name: "Incorrect Data",
			mock: func() {
				mock.ExpectExec("INSERT INTO locations").
					WithArgs(1, "Machu Picchu", "").
					WillReturnError(apperrors.ErrIncorrectQuery)
			},
			input: model.Location{
				LocationId: 1,
				Place:      "Machu Picchu",
				Country:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := repository.Insert(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
