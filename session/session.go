package session

import "github.com/gorilla/sessions"

// GameSession ...
type GameSession struct {
	CurrentSession sessions.Store
	SessionName    string
}

// NewSessionGame ...
func NewSessionGame(key string) *GameSession {
	return &GameSession{
		SessionName:    "hw",
		CurrentSession: sessions.NewCookieStore([]byte(key)),
	}
}
