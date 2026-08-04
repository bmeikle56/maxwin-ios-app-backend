package models

import "time"

// GameType matches iOS GameType raw values.
type GameType string

const (
	GameTypeCash       GameType = "Cash"
	GameTypeTournament GameType = "Tournament"
)

// User matches the iOS User model.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// StreetAction is a street-by-street action within a hand.
type StreetAction struct {
	ID       string   `json:"id"`
	Street   string   `json:"street"`
	Action   string   `json:"action"`
	PotAfter *float64 `json:"potAfter,omitempty"`
}

// HandDetail is the optional board / street breakdown for a hand.
type HandDetail struct {
	Board       *string        `json:"board,omitempty"`
	PotSize     *float64       `json:"potSize,omitempty"`
	Opponents   *int           `json:"opponents,omitempty"`
	VillainHand *string        `json:"villainHand,omitempty"`
	AllInStreet *string        `json:"allInStreet,omitempty"`
	Streets     []StreetAction `json:"streets,omitempty"`
}

// Hand is a single recorded poker hand within a session.
type Hand struct {
	ID         string      `json:"id"`
	HandNumber int         `json:"handNumber"`
	Position   string      `json:"position"`
	HoleCards  string      `json:"holeCards"`
	Result     float64     `json:"result"`
	Notes      *string     `json:"notes,omitempty"`
	Detail     *HandDetail `json:"detail,omitempty"`
}

// PokerSession matches the iOS PokerSession model.
type PokerSession struct {
	ID              string    `json:"id"`
	Date            time.Time `json:"date"`
	Venue           string    `json:"venue"`
	GameType        GameType  `json:"gameType"`
	Stakes          string    `json:"stakes"`
	DurationMinutes int       `json:"durationMinutes"`
	BuyIn           float64   `json:"buyIn"`
	CashOut         float64   `json:"cashOut"`
	Hands           []Hand    `json:"hands"`
}

// Profit is cashOut - buyIn.
func (s PokerSession) Profit() float64 {
	return s.CashOut - s.BuyIn
}

// EarningsDataPoint is a chart point for cumulative earnings.
type EarningsDataPoint struct {
	ID               string    `json:"id"`
	Date             time.Time `json:"date"`
	CumulativeProfit float64   `json:"cumulativeProfit"`
	PeriodProfit     float64   `json:"periodProfit"`
}

// DateRangeFilter matches iOS DateRangeFilter cases.
type DateRangeFilter string

const (
	RangeAllTime   DateRangeFilter = "allTime"
	RangeLastYear  DateRangeFilter = "lastYear"
	RangeLastMonth DateRangeFilter = "lastMonth"
)

// AuthRequest is the sign-in body.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PasswordResetRequest is the password-reset body.
type PasswordResetRequest struct {
	Username string `json:"username"`
}

// ErrorResponse is a standard API error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}
