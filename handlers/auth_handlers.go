package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"maxwin/middleware"
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

	resp, err := h.service.SignIn(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, mock.ErrEmptyFields) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandlers) SignOut(c *gin.Context) {
	// Stateless JWTs are discarded client-side. Endpoint kept for iOS symmetry.
	c.JSON(http.StatusOK, gin.H{"response": "signed out"})
}

func (h *AuthHandlers) DeleteAccount(c *gin.Context) {
	username := middleware.Username(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	h.service.DeleteAccount(username)
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
