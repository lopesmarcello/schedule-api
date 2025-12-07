package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) BindRoutes() {
	api.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := api.Router.Group("/api/v1")

	v1.POST("/user", api.handleRegister)
}
