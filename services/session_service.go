package services

import (
	"maxwin/mock"
	"maxwin/models"
)

type SessionService struct {
	store *mock.Store
}

func NewSessionService(store *mock.Store) *SessionService {
	return &SessionService{store: store}
}

func (s *SessionService) FetchSessions() []models.PokerSession {
	return s.store.ListSessions()
}

func (s *SessionService) GetSession(id string) (models.PokerSession, error) {
	return s.store.GetSession(id)
}

func (s *SessionService) CreateSession(session models.PokerSession) (models.PokerSession, error) {
	return s.store.CreateSession(session)
}

func (s *SessionService) UpdateSession(session models.PokerSession) (models.PokerSession, error) {
	return s.store.UpdateSession(session)
}

func (s *SessionService) DeleteSession(id string) error {
	return s.store.DeleteSession(id)
}
