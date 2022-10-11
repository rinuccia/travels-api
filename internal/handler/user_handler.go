package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rinuccia/travels-api/internal/model"
	"github.com/rinuccia/travels-api/internal/service"
	"github.com/rinuccia/travels-api/pkg/apperrors"
	"net/http"
)

type userHandler struct {
	repo service.User
}

func newUserHandler(r service.User) *userHandler {
	return &userHandler{
		repo: r,
	}
}

// getUserById godoc
// @Summary Returns user based on given ID
// @Tags user
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} model.User
// @Failure 404 {object} errResponse
// @Router /user/{id} [get]
func (h *userHandler) getUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, newErrResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)
}

// createUser godoc
// @Summary Create user
// @Tags user
// @Accept json
// @Produce json
// @Param input body model.User true "User Info"
// @Success 200 {object} model.User
// @Failure 400 {object} errResponse
// @Router /user/new [post]
func (h *userHandler) createUser(c *gin.Context) {
	user := model.User{}
	err := c.BindJSON(&user)
	validationErr := validate.Struct(user)
	if err != nil || validationErr != nil {
		c.JSON(http.StatusBadRequest, newErrResponse("invalid input body"))
		return
	}

	user, err = h.repo.Create(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)
}

// updateUser godoc
// @Summary Update user based on given ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Param input body model.User true "User Info"
// @Success 204
// @Failure 400,404 {object} errResponse
// @Router /user/{id} [put]
func (h *userHandler) updateUser(c *gin.Context) {
	id := c.Param("id")
	user := model.User{}
	err := c.BindJSON(&user)
	validationErr := validate.Struct(user)
	if err != nil || validationErr != nil {
		c.JSON(http.StatusBadRequest, newErrResponse("invalid input body"))
		return
	}

	err = h.repo.Update(id, user)
	if errors.Is(err, apperrors.ErrIncorrectQuery) {
		c.JSON(http.StatusBadRequest, newErrResponse(err.Error()))
		return
	}
	if errors.Is(err, apperrors.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, newErrResponse(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
