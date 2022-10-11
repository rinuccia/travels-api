package postgres

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestVisitRepo_FindAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newVisitRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		userId  string
		want    model.UserVisits
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				row := sqlmock.NewRows([]string{"user_id"}).AddRow(1)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs("1").WillReturnRows(row)
				rows := sqlmock.NewRows([]string{"place", "country", "visited_at", "mark"}).
					AddRow("Red Square", "RF", "2015-06-23", 4).
					AddRow("Grand Canyon", "USA", "2019-04-30", 5)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("1").WillReturnRows(rows)
			},
			userId: "1",
			want: model.UserVisits{
				Visits: []model.UserVisit{
					{"Red Square", "RF", "2015-06-23", 4},
					{"Grand Canyon", "USA", "2019-04-30", 5},
				},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs("1").WillReturnError(apperrors.ErrRecordNotFound)
			},
			userId:  "1",
			wantErr: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock()

			got, err := repository.FindAll(tt.userId)

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

func TestVisitRepo_Insert(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newVisitRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		input   model.Visit
		want    model.Visit
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO visits (.+)").
					WithArgs(1, 1, 2, "2019-06-15", 4).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: model.Visit{
				VisitId:    1,
				LocationId: 1,
				UserId:     2,
				VisitedAt:  "2019-06-15",
				Mark:       4,
			},
			want: model.Visit{
				VisitId:    1,
				LocationId: 1,
				UserId:     2,
				VisitedAt:  "2019-06-15",
				Mark:       4,
			},
		},
		{
			name: "Incorrect Data",
			mock: func() {
				mock.ExpectExec("INSERT INTO visits (.+)").
					WithArgs(1, 1, 2, "2019-06-15", 10).WillReturnError(apperrors.ErrIncorrectQuery)
			},
			input: model.Visit{
				VisitId:    1,
				LocationId: 1,
				UserId:     2,
				VisitedAt:  "2019-06-15",
				Mark:       10,
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

func TestVisitRepo_DeleteById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	repository := newVisitRepo(db)

	testTable := []struct {
		name    string
		mock    func()
		id      string
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM visits WHERE (.+)").WithArgs("1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: "1",
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM visits WHERE (.+)").WithArgs("1").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			id:      "1",
			wantErr: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock()

			err = repository.DeleteById(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
