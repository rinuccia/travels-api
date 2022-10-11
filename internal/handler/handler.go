package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rinuccia/travels-api/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rinuccia/travels-api/docs"
)

const (
	userURL      = "/user"
	locationURL  = "/location"
	locationsURL = "/locations"
	visitURL     = "/visit"
	visitsURL    = "/visits"
)

var validate = validator.New()

type Handler struct {
	*userHandler
	*locationHandler
	*visitHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		newUserHandler(service.User),
		newLocationHandler(service.Location),
		newVisitHandler(service.Visit),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET(userURL+"/:id", h.getUserById)
	router.POST(userURL+"/new", h.createUser)
	router.PUT(userURL+"/:id", h.updateUser)
	router.GET(locationsURL, h.getAllLocations)
	router.GET(locationURL+"/:id", h.getLocationById)
	router.GET(locationURL+"/:id/avg", h.getAvgRating)
	router.POST(locationURL+"/new", h.createLocation)
	router.GET(visitsURL+"/user/:id", h.getAllVisits)
	router.POST(visitURL+"/new", h.createVisit)
	router.DELETE(visitURL+"/:id", h.deleteVisitById)
}
