package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"maxwin/models"
	"maxwin/services"
)

type EarningsHandlers struct {
	service *services.EarningsService
}

func NewEarningsHandlers(service *services.EarningsService) *EarningsHandlers {
	return &EarningsHandlers{service: service}
}

func (h *EarningsHandlers) Fetch(c *gin.Context) {
	rangeParam := c.DefaultQuery("range", string(models.RangeAllTime))
	var filter models.DateRangeFilter

	switch rangeParam {
	case string(models.RangeAllTime), "All time":
		filter = models.RangeAllTime
	case string(models.RangeLastYear), "Last year":
		filter = models.RangeLastYear
	case string(models.RangeLastMonth), "Last month":
		filter = models.RangeLastMonth
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "range must be allTime, lastYear, or lastMonth",
		})
		return
	}

	c.JSON(http.StatusOK, h.service.FetchEarnings(filter))
}
