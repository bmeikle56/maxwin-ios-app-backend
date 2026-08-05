package services

import (
	"maxwin/auth"
	"maxwin/mock"
	"maxwin/models"
)

type AuthService struct {
	store *mock.Store
}

func NewAuthService(store *mock.Store) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) SignIn(username, password string) (models.AuthResponse, error) {
	user, err := s.store.SignIn(username, password)
	if err != nil {
		return models.AuthResponse{}, err
	}

	token, expiresAt, err := auth.MintToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) DeleteAccount(username string) {
	s.store.DeleteAccount(username)
}

func (s *AuthService) RequestPasswordReset(username string) error {
	return s.store.RequestPasswordReset(username)
}
