package auth

import (
	"context"
	"net/http"

	"github.com/yan0ks/codename-timecode/internal/entity"
)

type contextKey string

const userContextKey contextKey = "user"

func CurrentUser(ctx context.Context) (entity.Principal, bool) {
	p, ok := ctx.Value(userContextKey).(entity.Principal)
	return p, ok
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			principal, err := s.GetUserBySessionID(cookie.Value)
			if err == nil {
				ctx := context.WithValue(r.Context(), userContextKey, principal)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}
