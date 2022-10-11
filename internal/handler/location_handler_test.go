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

func TestLocationHandler_getAllLocations(t *testing.T) {
	type mockBehavior func(s *mock_service.MockLocation)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockLocation) {
				s.EXPECT().GetAll().Return(model.Locations{
					List: []model.Location{
						{1, "Red Square", "RF"},
						{2, "Eiffel Tower", "France"},
						{3, "Grand Canyon", "USA"},
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"list":[{"location_id":1,"place":"Red Square","country":"RF"},{"location_id":2,"place":"Eiffel Tower","country":"France"},{"location_id":3,"place":"Grand Canyon","country":"USA"}]}`,
		},
		{
			name: "Service Error",
			mockBehavior: func(s *mock_service.MockLocation) {
				s.EXPECT().GetAll().Return(model.Locations{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			location := mock_service.NewMockLocation(controller)
			test.mockBehavior(location)

			serv := &service.Service{Location: location}
			handle := NewHandler(serv)

			router := gin.New()
			router.GET("/locations", handle.getAllLocations)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/locations", nil)

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestLocationHandler_getLocationById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockLocation, id string)

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
			mockBehavior: func(s *mock_service.MockLocation, id string) {
				s.EXPECT().GetById(id).Return(model.Location{
					LocationId: 1,
					Place:      "Red Square",
					Country:    "RF"}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"location_id":1,"place":"Red Square","country":"RF"}`,
		},
		{
			name: "Not Found",
			id:   "1",
			mockBehavior: func(s *mock_service.MockLocation, id string) {
				s.EXPECT().GetById(id).Return(model.Location{}, apperrors.ErrRecordNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"record not found"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			location := mock_service.NewMockLocation(controller)
			test.mockBehavior(location, test.id)

			serv := &service.Service{Location: location}
			handle := NewHandler(serv)

			router := gin.New()
			router.GET("/location/:id", handle.getLocationById)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/location/1", nil)

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestLocationHandler_getAvgRating(t *testing.T) {
	type mockBehavior func(s *mock_service.MockLocation, id string)

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
			mockBehavior: func(s *mock_service.MockLocation, id string) {
				s.EXPECT().GetRating(id).Return(float32(4.5), nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"avg":4.5}`,
		},
		{
			name: "Not Found",
			id:   "1",
			mockBehavior: func(s *mock_service.MockLocation, id string) {
				s.EXPECT().GetRating(id).Return(float32(0), apperrors.ErrRecordNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"record not found"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			location := mock_service.NewMockLocation(controller)
			test.mockBehavior(location, test.id)

			serv := &service.Service{Location: location}
			handle := NewHandler(serv)

			router := gin.New()
			router.GET("/location/:id/avg", handle.getAvgRating)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/location/1/avg", nil)

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}

func TestLocationHandler_createLocation(t *testing.T) {
	type mockBehavior func(s *mock_service.MockLocation, location model.Location)

	testTable := []struct {
		name                 string
		inputBody            string
		inputLocation        model.Location
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"location_id":1,"place":"Machu Picchu","country":"Peru"}`,
			inputLocation: model.Location{
				LocationId: 1,
				Place:      "Machu Picchu",
				Country:    "Peru",
			},
			mockBehavior: func(s *mock_service.MockLocation, location model.Location) {
				s.EXPECT().Create(location).Return(location, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"location_id":1,"place":"Machu Picchu","country":"Peru"}`,
		},
		{
			name:      "Location exist",
			inputBody: `{"location_id":1,"place":"Machu Picchu","country":"Peru"}`,
			inputLocation: model.Location{
				LocationId: 1,
				Place:      "Machu Picchu",
				Country:    "Peru",
			},
			mockBehavior: func(s *mock_service.MockLocation, location model.Location) {
				s.EXPECT().Create(location).Return(model.Location{}, apperrors.ErrIncorrectQuery)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"incorrect query"}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			location := mock_service.NewMockLocation(controller)
			test.mockBehavior(location, test.inputLocation)

			serv := &service.Service{Location: location}
			handle := NewHandler(serv)

			router := gin.New()
			router.POST("/location/new", handle.createLocation)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/location/new", strings.NewReader(test.inputBody))

			router.ServeHTTP(w, r)

			body := strings.Trim(w.Body.String(), "\n")

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, body, test.expectedResponseBody)
		})
	}
}
