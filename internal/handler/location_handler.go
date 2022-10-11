package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/service"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"net/http"
)

type locationHandler struct {
	repo service.Location
}

func newLocationHandler(repository service.Location) *locationHandler {
	return &locationHandler{
		repo: repository,
	}
}

// getAllLocations godoc
// @Summary Returns a list of all locations
// @Tags location
// @Produce json
// @Success 200 {object} model.Locations
// @Failure 500 {object} errResponse
// @Router /locations [get]
func (h *locationHandler) getAllLocations(c *gin.Context) {
	locations, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrResponse("something went wrong"))
		return
	}

	c.JSON(http.StatusOK, locations)
}

// getLocationById godoc
// @Summary Returns location based on given ID
// @Tags location
// @Produce json
// @Param id path integer true "Location ID"
// @Success 200 {object} model.Location
// @Failure 404 {string} string
// @Router /location/{id} [get]
func (h *locationHandler) getLocationById(c *gin.Context) {
	id := c.Param("id")
	location, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, location)
}

// getAvgRating godoc
// @Summary Retrieves the average location rating based on given id
// @Tags location
// @Produce json
// @Param id path integer true "Location ID"
// @Success 200 {object} model.AvgRating
// @Failure 404 {object} errResponse
// @Router /location/{id}/avg [get]
func (h *locationHandler) getAvgRating(c *gin.Context) {
	id := c.Param("id")
	avg, err := h.repo.GetRating(id)
	if errors.Is(err, apperrors.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, newErrResponse(err.Error()))
		return
	}

	rating := model.AvgRating{Avg: avg}
	c.JSON(http.StatusOK, rating)
}

// createLocation godoc
// @Summary Create Location
// @Tags location
// @Accept json
// @Produce json
// @Param input body model.Location true "Location Info"
// @Success 200 {object} model.Location
// @Failure 400 {object} errResponse
// @Router /location/new [post]
func (h *locationHandler) createLocation(c *gin.Context) {
	location := model.Location{}
	err := c.BindJSON(&location)
	validationErr := validate.Struct(location)
	if err != nil || validationErr != nil {
		c.JSON(http.StatusBadRequest, newErrResponse("invalid input body"))
		return
	}

	location, err = h.repo.Create(location)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, location)
}
