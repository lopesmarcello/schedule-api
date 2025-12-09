package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lopesmarcello/schedule-api/internal/services"
)

type API struct {
	Router              *gin.Engine
	UserService         services.UserService
	AvailabilityService services.AvailabilityService
	AppointmentsService services.AppointmentsService
}

func NewAPI(userService services.UserService, availabilityService services.AvailabilityService, appointmentsService services.AppointmentsService) *API {
	router := gin.Default()
	api := &API{
		Router:              router,
		UserService:         userService,
		AvailabilityService: availabilityService,
		AppointmentsService: appointmentsService,
	}

	router.Use(ErrorHandler())

	api.BindRoutes()

	return api
}
