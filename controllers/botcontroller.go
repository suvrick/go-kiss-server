package controllers

import (
	"net/http"

	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/store"
	"github.com/suvrick/go-kiss-server/until"

	"github.com/gorilla/mux"
)

// BotController ...
type BotController struct {
	router        *mux.Router
	botRepository *store.BotRepository
	session       *session.GameSession
}

// NewBotController ...
func NewBotController(router *mux.Router, botRepository *store.BotRepository, sg *session.GameSession) *BotController {
	ctrl := &BotController{
		router:        router,
		botRepository: botRepository,
		session:       sg,
	}

	bot := ctrl.router.PathPrefix("/bot").Subrouter()
	bot.HandleFunc("/all", ctrl.all()).Methods("GET")

	return ctrl
}

// GET all
func (ctrl *BotController) all() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(middlewares.CtxKeyUser)

		until.WriteResponse(w, r, 200, map[string]interface{}{
			"result": "ok",
			"user":   user,
		}, nil)
	}
}
