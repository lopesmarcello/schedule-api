package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (api *API) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewHTTPError(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := api.UserService.CreateUser(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password)
	if err != nil {
		c.Error(NewHTTPError(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (api *API) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewHTTPError(http.StatusBadRequest, err.Error()))
		return
	}
	user, err := api.UserService.AuthenticateUser(
		c.Request.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		c.Error(NewHTTPError(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, user)
}
