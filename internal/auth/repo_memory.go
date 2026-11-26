package auth

import (
	"errors"
	"sync"

	"github.com/yan0ks/codename-timecode/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type MemoryRepo struct {
	mu       sync.RWMutex
	users    map[string]entity.User
	sessions map[string]entity.Session
}

func NewMemoryRepo() *MemoryRepo {
	r := &MemoryRepo{
		users:    make(map[string]entity.User),
		sessions: make(map[string]entity.Session),
	}
	r.seedDemoUsers()
	return r
}

func (r *MemoryRepo) seedDemoUsers() {
	viewerHash, _ := bcrypt.GenerateFromPassword([]byte("viewer"), bcrypt.DefaultCost)
	editorHash, _ := bcrypt.GenerateFromPassword([]byte("editor"), bcrypt.DefaultCost)

	r.users["viewer"] = entity.User{ID: "u1", Login: "viewer", PasswordHash: string(viewerHash), Role: entity.RoleViewer}
	r.users["editor"] = entity.User{ID: "u2", Login: "editor", PasswordHash: string(editorHash), Role: entity.RoleEditor}
}

func (r *MemoryRepo) GetUserByLogin(login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[login]
	if !ok {
		return entity.User{}, errors.New("not found")
	}
	return user, nil
}

func (r *MemoryRepo) GetUserByID(id string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return entity.User{}, errors.New("not found")
}

func (r *MemoryRepo) CreateSession(session entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *MemoryRepo) GetSessionByID(id string) (entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[id]
	if !ok {
		return entity.Session{}, errors.New("not found")
	}
	return session, nil
}

func (r *MemoryRepo) DeleteSession(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}
