package entity

import (
	"errors"
	"time"
)

type Role string

const (
	RoleViewer Role = "viewer"
	RoleEditor Role = "editor"
)

type User struct {
	ID           string
	Login        string
	PasswordHash string
	Role         Role
}

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type Principal struct {
	ID   string
	Role Role
}

var (
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUnauthorized       = errors.New("unauthorized")
)
