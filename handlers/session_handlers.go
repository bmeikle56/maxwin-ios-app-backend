package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"maxwin/mock"
	"maxwin/models"
	"maxwin/services"
)

type SessionHandlers struct {
	service *services.SessionService
}

func NewSessionHandlers(service *services.SessionService) *SessionHandlers {
	return &SessionHandlers{service: service}
}

func (h *SessionHandlers) List(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.FetchSessions())
}

func (h *SessionHandlers) Get(c *gin.Context) {
	session, err := h.service.GetSession(c.Param("id"))
	if err != nil {
		if errors.Is(err, mock.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandlers) Create(c *gin.Context) {
	var session models.PokerSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	created, err := h.service.CreateSession(session)
	if err != nil {
		if errors.Is(err, mock.ErrInvalidSession) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *SessionHandlers) Update(c *gin.Context) {
	var session models.PokerSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid JSON"})
		return
	}

	session.ID = c.Param("id")
	updated, err := h.service.UpdateSession(session)
	if err != nil {
		switch {
		case errors.Is(err, mock.ErrNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		case errors.Is(err, mock.ErrInvalidSession):
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *SessionHandlers) Delete(c *gin.Context) {
	if err := h.service.DeleteSession(c.Param("id")); err != nil {
		if errors.Is(err, mock.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
