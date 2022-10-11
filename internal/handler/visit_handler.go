package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/service"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"net/http"
)

type visitHandler struct {
	repo service.Visit
}

func newVisitHandler(repository service.Visit) *visitHandler {
	return &visitHandler{
		repo: repository,
	}
}

// getAllVisits godoc
// @Summary Returns a list of all user visits
// @Tags visit
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} model.UserVisits
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /visits/user/{id} [get]
func (h *visitHandler) getAllVisits(c *gin.Context) {
	id := c.Param("id")
	visits, err := h.repo.GetAll(id)
	if errors.Is(err, apperrors.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, newErrResponse(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrResponse("something went wrong"))
		return
	}

	c.JSON(http.StatusOK, visits)
}

// createVisit godoc
// @Summary Create Visit
// @Tags visit
// @Accept json
// @Produce json
// @Param input body model.Visit true "Visit Info"
// @Success 200 {object} model.Visit
// @Failure 400 {object} errResponse
// @Router /visit/new [post]
func (h *visitHandler) createVisit(c *gin.Context) {
	visit := model.Visit{}
	err := c.BindJSON(&visit)
	validationErr := validate.Struct(visit)
	if err != nil || validationErr != nil {
		c.JSON(http.StatusBadRequest, newErrResponse("invalid input body"))
		return
	}

	visit, err = h.repo.Create(visit)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, visit)
}

// deleteVisitById godoc
// @Summary Removes visit based on given ID
// @Tags visit
// @Produce json
// @Param id path integer true "Visit ID"
// @Success 204
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /visit/{id} [delete]
func (h *visitHandler) deleteVisitById(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.DeleteById(id)
	if errors.Is(err, apperrors.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, newErrResponse(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrResponse("something went wrong"))
		return
	}

	c.Status(http.StatusNoContent)
}
