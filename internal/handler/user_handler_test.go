package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/service"
	mock_service "github.com/rinuccia/travels-api/internal/service/mocks"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserHandler_getUserById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, id string)

	testTable := []struct {
		name                 string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id:   "1",
			mockBehavior: func(s *mock_service.MockUser, id string) {
				s.EXPECT().GetById(id).Return(model.User{
					UserId:    1,
					Email:     "test@gmail.com",
					FirstName: "John",
					LastName:  "Smith",
					Gender:    "m",
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":1,"email":"test@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
		},
		{
			name: "Not Found",
			id:   "1",
			mockBehavior: func(s *mock_service.MockUser, id string) {
				s.EXPECT().GetById(id).Return(model.User{}, apperrors.ErrRecordNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"record not found"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			user := mock_service.NewMockUser(controller)
			test.mockBehavior(user, test.id)

			serv := &service.Service{User: user}
			handle := NewHandler(serv)

			router := gin.New()
			router.GET("/user/:id", handle.getUserById)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/user/1", strings.NewReader(test.id))

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestUserHandler_createUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user model.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"user_id":1,"email":"test@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
			inputUser: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			mockBehavior: func(s *mock_service.MockUser, user model.User) {
				s.EXPECT().Create(user).Return(user, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"user_id":1,"email":"test@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
		},
		{
			name:      "User exists",
			inputBody: `{"user_id":1,"email":"smith@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
			inputUser: model.User{
				UserId:    1,
				Email:     "smith@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			mockBehavior: func(s *mock_service.MockUser, user model.User) {
				s.EXPECT().Create(user).Return(model.User{}, apperrors.ErrIncorrectQuery)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect query"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			user := mock_service.NewMockUser(controller)
			test.mockBehavior(user, test.inputUser)

			serv := &service.Service{User: user}
			handle := NewHandler(serv)

			router := gin.New()
			router.POST("/user/new", handle.createUser)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/user/new", strings.NewReader(test.inputBody))

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestUserHandler_updateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user model.User, id string)

	testTable := []struct {
		name                 string
		id                   string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			id:        "1",
			inputBody: `{"user_id":1,"email":"test@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
			inputUser: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			mockBehavior: func(s *mock_service.MockUser, user model.User, id string) {
				s.EXPECT().Update(id, user).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:      "Not Found",
			id:        "1",
			inputBody: `{"user_id":1,"email":"test@gmail.com","first_name":"John","last_name":"Smith","gender":"m"}`,
			inputUser: model.User{
				UserId:    1,
				Email:     "test@gmail.com",
				FirstName: "John",
				LastName:  "Smith",
				Gender:    "m",
			},
			mockBehavior: func(s *mock_service.MockUser, user model.User, id string) {
				s.EXPECT().Update(id, user).Return(apperrors.ErrRecordNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			user := mock_service.NewMockUser(controller)
			test.mockBehavior(user, test.inputUser, test.id)

			serv := &service.Service{User: user}
			handle := NewHandler(serv)

			router := gin.New()
			router.PUT("/user/:id", handle.updateUser)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/user/1", strings.NewReader(test.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}
