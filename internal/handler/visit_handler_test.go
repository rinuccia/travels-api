package handler

import (
	"errors"
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

func TestVisitHandler_getAllVisits(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVisit, id string)

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
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().GetAll(id).Return(model.UserVisits{
					Visits: []model.UserVisit{
						{Place: "Eiffel Tower", Country: "France", VisitedAt: "2015-06-12", Mark: 4},
						{Place: "Grand Canyon", Country: "USA", VisitedAt: "2019-09-02", Mark: 3},
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"visits":[{"place":"Eiffel Tower","country":"France","visited_at":"2015-06-12","mark":4},{"place":"Grand Canyon","country":"USA","visited_at":"2019-09-02","mark":3}]}`,
		},
		{
			name: "Not Found",
			id:   "1",
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().GetAll(id).Return(model.UserVisits{}, apperrors.ErrRecordNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"record not found"}`,
		},
		{
			name: "Service Error",
			id:   "1",
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().GetAll(id).Return(model.UserVisits{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			visit := mock_service.NewMockVisit(controller)
			test.mockBehavior(visit, test.id)

			serv := &service.Service{Visit: visit}
			handle := NewHandler(serv)

			router := gin.New()
			router.GET("/visits/user/:id", handle.getAllVisits)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/visits/user/1", nil)

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestVisitHandler_createVisit(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVisit, visit model.Visit)

	testTable := []struct {
		name                 string
		inputBody            string
		inputVisit           model.Visit
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"visit_id":1,"location_id":1,"user_id":1,"visited_at":"2018-10-16","mark":3}`,
			inputVisit: model.Visit{
				VisitId:    1,
				LocationId: 1,
				UserId:     1,
				VisitedAt:  "2018-10-16",
				Mark:       3,
			},
			mockBehavior: func(s *mock_service.MockVisit, visit model.Visit) {
				s.EXPECT().Create(visit).Return(visit, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"visit_id":1,"location_id":1,"user_id":1,"visited_at":"2018-10-16","mark":3}`,
		},
		{
			name:      "Visit exist",
			inputBody: `{"visit_id":1,"location_id":1,"user_id":1,"visited_at":"2018-10-16","mark":5}`,
			inputVisit: model.Visit{
				VisitId:    1,
				LocationId: 1,
				UserId:     1,
				VisitedAt:  "2018-10-16",
				Mark:       5,
			},
			mockBehavior: func(s *mock_service.MockVisit, visit model.Visit) {
				s.EXPECT().Create(visit).Return(visit, apperrors.ErrIncorrectQuery)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect query"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			visit := mock_service.NewMockVisit(controller)
			test.mockBehavior(visit, test.inputVisit)

			serv := &service.Service{Visit: visit}
			handle := NewHandler(serv)

			router := gin.New()
			router.POST("/visit/new", handle.createVisit)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/visit/new", strings.NewReader(test.inputBody))

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestVisitHandler_deleteVisitById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockVisit, id string)

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
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().DeleteById(id).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Not Found",
			id:   "1",
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().DeleteById(id).Return(apperrors.ErrRecordNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"record not found"}`,
		},
		{
			name: "Service Error",
			id:   "1",
			mockBehavior: func(s *mock_service.MockVisit, id string) {
				s.EXPECT().DeleteById(id).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			visit := mock_service.NewMockVisit(controller)
			test.mockBehavior(visit, test.id)

			serv := &service.Service{Visit: visit}
			handle := NewHandler(serv)

			router := gin.New()
			router.DELETE("/visit/:id", handle.deleteVisitById)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/visit/1", nil)

			router.ServeHTTP(w, r)

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}
