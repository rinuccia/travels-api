package postgres

import (
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestUserRepo_FindById(t *testing.T) {
	mockDB, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer mockDB.Close()

	repository := newUserRepo(mockDB)

	testTable := []struct {
		name    string
		mock    func()
		id      string
		want    model.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := mock.NewRows([]string{"user_id", "email", "first_name", "last_name", "gender"}).
					AddRow("1", "test@gmail.com", "John", "Smith", "m")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("1").WillReturnRows(rows)
			},
			id: "1",
			want: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
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

func TestUserRepo_Insert(t *testing.T) {
	mockDB, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer mockDB.Close()

	repository := newUserRepo(mockDB)

	testTable := []struct {
		name    string
		mock    func()
		input   model.User
		want    model.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(1, "test@gmail.com", "John", "Smith", "m").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			want: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
		},
		{
			name: "Incorrect Data",
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(1, "test@gmail.com", "John", "Smith", "").
					WillReturnError(apperrors.ErrIncorrectQuery)
			},
			input: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "",
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

func TestUserRepo_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer mockDB.Close()

	repository := newUserRepo(mockDB)

	testTable := []struct {
		name            string
		mock            func()
		id              string
		input           model.User
		wantErr         bool
		expectedErrType error
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("UPDATE users SET (.+) WHERE (.+)").
					WithArgs("test@gmail.com", "John", "Smith", "m", "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id: "1",
			input: model.User{
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
		},
		{
			name: "Incorrect Data",
			mock: func() {
				mock.ExpectExec("UPDATE users SET (.+) WHERE (.+)").
					WithArgs("test@gmail.com", "John", "Smith", "", "1").
					WillReturnError(apperrors.ErrIncorrectQuery)
			},
			id: "1",
			input: model.User{
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "",
			},
			wantErr:         true,
			expectedErrType: apperrors.ErrIncorrectQuery,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("UPDATE users SET (.+) WHERE (.+)").
					WithArgs("test@gmail.com", "John", "Smith", "m", "1").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			id: "1",
			input: model.User{
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			wantErr:         true,
			expectedErrType: apperrors.ErrRecordNotFound,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err = repository.Update(tt.id, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrType != nil {
					assert.Equal(t, tt.expectedErrType, err)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
