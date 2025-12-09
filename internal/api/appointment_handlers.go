package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) handleGetSlots(c *gin.Context) {
	userSlug := c.Param("slug") // Assume route /:slug/slots/:date
	dateStr := c.Param("date")

	// Find userID by slug (assume a service method)
	user, err := api.UserService.GetUserBySlug(c.Request.Context(), userSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	slots, _, err := api.AppointmentsService.GetAvailableSpots(c.Request.Context(), user.ID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if slots == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}

	slotStrings := make([]string, len(slots))
	for i, slot := range slots {
		slotStrings[i] = slot.Start[:5]
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": slotStrings})
}
