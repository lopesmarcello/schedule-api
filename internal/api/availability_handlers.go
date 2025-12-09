package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Availability struct {
	DayOfWeek int    `json:"day_of_week" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type SetAvailabilityReq struct {
	UserID       int            `json:"user_id" binding:"required"`
	Availability []Availability `json:"availability" binding:"required"`
}

func (api *API) handleSetAvailability(c *gin.Context) {
	var req SetAvailabilityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, avail := range req.Availability {
		err := api.AvailabilityService.SetAvailability(
			c.Request.Context(),
			req.UserID,
			avail.DayOfWeek,
			avail.StartTime,
			avail.EndTime)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": "Availability set"})
}
