package services

import (
	"maxwin/mock"
	"maxwin/models"
)

type AuthService struct {
	store *mock.Store
}

func NewAuthService(store *mock.Store) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) SignIn(username, password string) (models.User, error) {
	return s.store.SignIn(username, password)
}

func (s *AuthService) DeleteAccount(username string) {
	s.store.DeleteAccount(username)
}

func (s *AuthService) RequestPasswordReset(username string) error {
	return s.store.RequestPasswordReset(username)
}
