package mock

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"maxwin/models"
)

var (
	ErrNotFound       = errors.New("session not found")
	ErrInvalidSession = errors.New("enter a venue and stakes to save this session")
	ErrEmptyFields    = errors.New("username and password are required")
)

// Store is an in-memory stand-in for Postgres until the real DB is wired up.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]models.PokerSession
	users    map[string]models.User // username -> user
}

func NewStore() *Store {
	s := &Store{
		sessions: make(map[string]models.PokerSession),
		users:    make(map[string]models.User),
	}
	for _, session := range seedSessions() {
		s.sessions[session.ID] = session
	}
	return s
}

func (s *Store) ListSessions() []models.PokerSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.PokerSession, 0, len(s.sessions))
	for _, session := range s.sessions {
		out = append(out, session)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Date.After(out[j].Date)
	})
	return out
}

func (s *Store) GetSession(id string) (models.PokerSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return models.PokerSession{}, ErrNotFound
	}
	return session, nil
}

func (s *Store) CreateSession(session models.PokerSession) (models.PokerSession, error) {
	if err := validateSession(session); err != nil {
		return models.PokerSession{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if session.ID == "" {
		session.ID = uuid.NewString()
	}
	if session.Hands == nil {
		session.Hands = []models.Hand{}
	}
	s.sessions[session.ID] = session
	return session, nil
}

func (s *Store) UpdateSession(session models.PokerSession) (models.PokerSession, error) {
	if err := validateSession(session); err != nil {
		return models.PokerSession{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[session.ID]; !ok {
		return models.PokerSession{}, ErrNotFound
	}
	if session.Hands == nil {
		session.Hands = []models.Hand{}
	}
	s.sessions[session.ID] = session
	return session, nil
}

func (s *Store) DeleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[id]; !ok {
		return ErrNotFound
	}
	delete(s.sessions, id)
	return nil
}

func (s *Store) SignIn(username, password string) (models.User, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return models.User{}, ErrEmptyFields
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if user, ok := s.users[username]; ok {
		return user, nil
	}

	user := models.User{
		ID:       uuid.NewString(),
		Username: username,
	}
	s.users[username] = user
	return user, nil
}

func (s *Store) DeleteAccount(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, strings.TrimSpace(username))
}

func (s *Store) RequestPasswordReset(username string) error {
	if strings.TrimSpace(username) == "" {
		return ErrEmptyFields
	}
	return nil
}

func validateSession(session models.PokerSession) error {
	if strings.TrimSpace(session.Venue) == "" || strings.TrimSpace(session.Stakes) == "" {
		return ErrInvalidSession
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

func daysAgo(days, hour int) time.Time {
	now := time.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	return base.AddDate(0, 0, -days)
}

func seedSessions() []models.PokerSession {
	return []models.PokerSession{
		{
			ID:              uuid.NewString(),
			Date:            daysAgo(2, 19),
			Venue:           "Bellagio",
			GameType:        models.GameTypeCash,
			Stakes:          "2/5 NL",
			DurationMinutes: 210,
			BuyIn:           500,
			CashOut:         1280,
			Hands: []models.Hand{
				{
					ID: uuid.NewString(), HandNumber: 1, Position: "BTN", HoleCards: "A♠ K♠", Result: 340,
					Notes: ptr("Set over set on the turn"),
					Detail: &models.HandDetail{
						Board: ptr("A♥ K♦ 7♣ 2♠ 9♥"), PotSize: ptr(780.0), Opponents: ptr(1),
						VillainHand: ptr("7♠ 7♦"), AllInStreet: ptr("Turn"),
						Streets: []models.StreetAction{
							{ID: uuid.NewString(), Street: "Preflop", Action: "Raise to $25, called", PotAfter: ptr(55.0)},
							{ID: uuid.NewString(), Street: "Flop", Action: "Bet $40, called", PotAfter: ptr(135.0)},
							{ID: uuid.NewString(), Street: "Turn", Action: "Villain shoves, snap call", PotAfter: ptr(780.0)},
						},
					},
				},
				{
					ID: uuid.NewString(), HandNumber: 2, Position: "CO", HoleCards: "Q♥ Q♦", Result: -85,
					Notes: ptr("Folded to river jam"),
					Detail: &models.HandDetail{
						Board: ptr("J♣ 8♦ 2♥ T♠ A♣"), PotSize: ptr(210.0), Opponents: ptr(1),
						Streets: []models.StreetAction{
							{ID: uuid.NewString(), Street: "River", Action: "Faced jam, folded", PotAfter: ptr(210.0)},
						},
					},
				},
				{ID: uuid.NewString(), HandNumber: 3, Position: "SB", HoleCards: "J♣ T♣", Result: 210, Notes: ptr("Flush vs top pair")},
				{ID: uuid.NewString(), HandNumber: 4, Position: "UTG", HoleCards: "A♦ Q♦", Result: -120},
			},
		},
		{
			ID:              uuid.NewString(),
			Date:            daysAgo(8, 14),
			Venue:           "Local home game",
			GameType:        models.GameTypeTournament,
			Stakes:          "$100 buy-in",
			DurationMinutes: 320,
			BuyIn:           100,
			CashOut:         0,
			Hands: []models.Hand{
				{
					ID: uuid.NewString(), HandNumber: 1, Position: "MP", HoleCards: "K♥ K♦", Result: -100,
					Notes: ptr("Busted with overpair"),
					Detail: &models.HandDetail{
						Board: ptr("K♣ 9♠ 4♦ Q♥ A♠"), PotSize: ptr(240.0), Opponents: ptr(1),
						VillainHand: ptr("A♦ Q♦"), AllInStreet: ptr("River"),
					},
				},
				{ID: uuid.NewString(), HandNumber: 2, Position: "BTN", HoleCards: "A♣ 9♣", Result: 45, Notes: ptr("Chip up early")},
			},
		},
		{
			ID: uuid.NewString(), Date: daysAgo(18, 19), Venue: "ARIA", GameType: models.GameTypeCash,
			Stakes: "1/3 NL", DurationMinutes: 180, BuyIn: 300, CashOut: 145,
			Hands: []models.Hand{
				{ID: uuid.NewString(), HandNumber: 1, Position: "BB", HoleCards: "7♥ 7♠", Result: -90, Notes: ptr("Coolered by overpair")},
				{ID: uuid.NewString(), HandNumber: 2, Position: "HJ", HoleCards: "A♠ J♥", Result: -65},
			},
		},
		{
			ID: uuid.NewString(), Date: daysAgo(40, 19), Venue: "Online - Ignition", GameType: models.GameTypeCash,
			Stakes: "25NL", DurationMinutes: 95, BuyIn: 100, CashOut: 246,
			Hands: []models.Hand{
				{ID: uuid.NewString(), HandNumber: 1, Position: "CO", HoleCards: "T♠ T♥", Result: 88},
				{
					ID: uuid.NewString(), HandNumber: 2, Position: "BTN", HoleCards: "A♥ 5♥", Result: 58,
					Notes: ptr("Rivered flush"),
					Detail: &models.HandDetail{
						Board: ptr("K♥ 9♥ 2♣ 3♦ 7♥"), PotSize: ptr(92.0), Opponents: ptr(2),
						Streets: []models.StreetAction{
							{ID: uuid.NewString(), Street: "River", Action: "Bet half pot, both fold", PotAfter: ptr(92.0)},
						},
					},
				},
			},
		},
		{
			ID: uuid.NewString(), Date: daysAgo(95, 19), Venue: "Commerce", GameType: models.GameTypeCash,
			Stakes: "5/10 NL", DurationMinutes: 260, BuyIn: 1500, CashOut: 920,
			Hands: []models.Hand{
				{ID: uuid.NewString(), HandNumber: 1, Position: "UTG", HoleCards: "A♣ A♦", Result: -420, Notes: ptr("Lost to rivered straight")},
				{ID: uuid.NewString(), HandNumber: 2, Position: "BTN", HoleCards: "K♠ Q♠", Result: -160},
			},
		},
		{
			ID: uuid.NewString(), Date: daysAgo(200, 19), Venue: "WSOP Satellite", GameType: models.GameTypeTournament,
			Stakes: "$250 buy-in", DurationMinutes: 410, BuyIn: 250, CashOut: 1850,
			Hands: []models.Hand{
				{
					ID: uuid.NewString(), HandNumber: 1, Position: "HJ", HoleCards: "A♠ K♦", Result: 620,
					Notes: ptr("Final table double"),
					Detail: &models.HandDetail{
						Board: ptr("A♣ 8♦ 4♠ K♥ 2♣"), PotSize: ptr(1240.0), Opponents: ptr(1),
						VillainHand: ptr("A♥ Q♠"), AllInStreet: ptr("Flop"),
						Streets: []models.StreetAction{
							{ID: uuid.NewString(), Street: "Flop", Action: "Shove called", PotAfter: ptr(1240.0)},
						},
					},
				},
				{ID: uuid.NewString(), HandNumber: 2, Position: "SB", HoleCards: "9♣ 9♦", Result: 310},
			},
		},
	}
}
