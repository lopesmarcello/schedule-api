package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lopesmarcello/schedule-api/internal/services"
)

type API struct {
	Router      *gin.Engine
	UserService services.UserService
}
