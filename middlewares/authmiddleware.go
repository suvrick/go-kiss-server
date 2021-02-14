package middlewares

import (
	"context"
	"fmt"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/store"
	"github.com/suvrick/go-kiss-server/until"
	"net/http"
	"strings"
)

type ctxKey int8

// CtxKeyUser ...
const CtxKeyUser ctxKey = iota

var AuthMiddlewareInstance *AuthMiddleware

// AuthMiddleware ...
type AuthMiddleware struct {
	s    *session.GameSession
	repo *store.UserRepository
}

// NewAuthMiddleWare ...
func NewAuthMiddleWare(s *session.GameSession, repo *store.UserRepository) {
	AuthMiddlewareInstance = &AuthMiddleware{
		s:    s,
		repo: repo,
	}
}

// Do ...
func (m *AuthMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Start middleware")

		if strings.Compare(r.URL.Path, "/user/login") == 0 ||
			strings.Compare(r.URL.Path, "/user/register") == 0 {
			next.ServeHTTP(w, r)
			return
		}

		session, err := m.s.CurrentSession.Get(r, m.s.SessionName)
		if err != nil {
			until.WriteResponse(w, r, 403, nil, errors.ErrNotAuthenticated)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			until.WriteResponse(w, r, 403, nil, errors.ErrNotAuthenticated)
			return
		}

		u, err := m.repo.Find(id.(int))

		fmt.Println(u)

		if err != nil {
			until.WriteResponse(w, r, 403, nil, errors.ErrNotAuthenticated)
			return
		}

		context := context.WithValue(r.Context(), CtxKeyUser, u)
		requestContext := r.WithContext(context)
		next.ServeHTTP(w, requestContext)
	})
}
