package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"maxwin/mock"
	"maxwin/models"
	"maxwin/services"
)

type AuthHandlers struct {
	service *services.AuthService
}

func NewAuthHandlers(service *services.AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (h *AuthHandlers) SignIn(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	user, err := h.service.SignIn(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, mock.ErrEmptyFields) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandlers) SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "signed out"})
}

func (h *AuthHandlers) DeleteAccount(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	h.service.DeleteAccount(req.Username)
	c.JSON(http.StatusOK, gin.H{"response": "account deleted"})
}

func (h *AuthHandlers) PasswordReset(c *gin.Context) {
	var req models.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	if err := h.service.RequestPasswordReset(req.Username); err != nil {
		if errors.Is(err, mock.ErrEmptyFields) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": "password reset requested"})
}
