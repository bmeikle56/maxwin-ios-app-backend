package services

import (
	"sort"
	"time"

	"maxwin/mock"
	"maxwin/models"
)

type EarningsService struct {
	store *mock.Store
}

func NewEarningsService(store *mock.Store) *EarningsService {
	return &EarningsService{store: store}
}

func (s *EarningsService) FetchEarnings(rangeFilter models.DateRangeFilter) []models.EarningsDataPoint {
	sessions := s.store.ListSessions()
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Date.Before(sessions[j].Date)
	})

	start := rangeStart(rangeFilter, time.Now())
	var running float64
	points := make([]models.EarningsDataPoint, 0, len(sessions))

	for _, session := range sessions {
		if start != nil && session.Date.Before(*start) {
			continue
		}
		profit := session.Profit()
		running += profit
		points = append(points, models.EarningsDataPoint{
			ID:               session.ID,
			Date:             session.Date,
			CumulativeProfit: running,
			PeriodProfit:     profit,
		})
	}
	return points
}

func rangeStart(filter models.DateRangeFilter, now time.Time) *time.Time {
	switch filter {
	case models.RangeLastYear:
		t := now.AddDate(-1, 0, 0)
		return &t
	case models.RangeLastMonth:
		t := now.AddDate(0, -1, 0)
		return &t
	default:
		return nil
	}
}
