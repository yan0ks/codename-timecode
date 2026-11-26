package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/yan0ks/codename-timecode/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetUserByLogin(login string) (entity.User, error)
	CreateSession(session entity.Session) error
	GetSessionByID(id string) (entity.Session, error)
	DeleteSession(id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(login, password string) (entity.Principal, string, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return entity.Principal{}, "", entity.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return entity.Principal{}, "", entity.ErrInvalidCredentials
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return entity.Principal{}, "", err
	}
	sessionID := hex.EncodeToString(b)

	session := entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.repo.CreateSession(session); err != nil {
		return entity.Principal{}, "", err
	}

	principal := entity.Principal{ID: user.ID, Role: user.Role}
	return principal, sessionID, nil
}

func (s *Service) GetUserBySessionID(sessionID string) (entity.Principal, error) {
	session, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return entity.Principal{}, entity.ErrUnauthorized
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repo.DeleteSession(sessionID)
		return entity.Principal{}, entity.ErrUnauthorized
	}

	return entity.Principal{ID: session.UserID, Role: entity.RoleViewer}, nil
}

func (s *Service) Logout(sessionID string) {
	_ = s.repo.DeleteSession(sessionID)
}
