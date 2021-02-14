package controllers

import (
	"errors"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/store"
	"github.com/suvrick/go-kiss-server/until"
	"net/http"

	"github.com/gorilla/mux"
)

// UserController ...
type UserController struct {
	router         *mux.Router
	userRepository *store.UserRepository
	session        *session.GameSession
}

// NewUserController ...
func NewUserController(router *mux.Router, userRepository *store.UserRepository, sg *session.GameSession) *UserController {
	ctrl := &UserController{
		router:         router,
		userRepository: userRepository,
		session:        sg,
	}

	user := ctrl.router.PathPrefix("/user").Subrouter()

	user.Use(middlewares.AuthMiddlewareInstance.Do)

	user.HandleFunc("/login", ctrl.loginHandler()).Methods("POST")
	user.HandleFunc("/register", ctrl.registerHandler()).Methods("POST")
	user.HandleFunc("/logout", ctrl.logoutHandler()).Methods("GET")
	user.HandleFunc("/who", ctrl.whoHandler()).Methods("GET")

	return ctrl
}

// GET who
func (ctrl *UserController) whoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(middlewares.CtxKeyUser)
		until.WriteResponse(w, r, 200, map[string]interface{}{
			"result": "ok",
			"user":   user,
		}, nil)
	}
}

// POST Register
func (ctrl *UserController) registerHandler() http.HandlerFunc {

	type register struct {
		Email    string `json:email`
		Password string `json:password`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reg := &register{}
		if err := until.JSONBind(r, reg); err != nil {
			until.WriteResponse(w, r, http.StatusBadRequest, nil, err)
			return
		}

		if len(reg.Email) < 5 || len(reg.Password) < 5 {
			until.WriteResponse(w, r, http.StatusBadRequest, nil, errors.New("Слишком короткий логин или пароль"))
			return
		}

		u := &model.User{
			Email:    reg.Email,
			Password: reg.Password,
		}

		id, err := ctrl.userRepository.Create(u)
		if err != nil {
			until.WriteResponse(w, r, http.StatusBadRequest, nil, errors.New("Ошибка при регистрации.Попробуйте другой e-mail"))
			return
		}

		u.Sanitize()
		until.WriteResponse(w, r, 200, map[string]interface{}{
			"result": "ok",
			"id":     id,
		}, nil)
	}
}

// POST Login
func (ctrl *UserController) loginHandler() http.HandlerFunc {

	type login struct {
		Email    string `json: email`
		Password string `json: password`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		log := &login{}
		if err := until.JSONBind(r, log); err != nil {
			until.WriteResponse(w, r, http.StatusBadRequest, nil, err)
			return
		}

		u, err := ctrl.userRepository.FindByEmail(log.Email)
		if err != nil || !u.ComparePassword(log.Password) {
			until.WriteResponse(w, r, 200, nil, errors.New("Не верный логин или пароль"))
			return
		}

		userSession, err := ctrl.session.CurrentSession.Get(r, ctrl.session.SessionName)
		if err != nil {
			until.WriteResponse(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		userSession.Values["user_id"] = u.ID
		if err := ctrl.session.CurrentSession.Save(r, w, userSession); err != nil {
			until.WriteResponse(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		until.WriteResponse(w, r, 200, map[string]interface{}{
			"result": "ok",
			"user":   u,
		}, nil)
	}
}

// GET Logout
func (ctrl *UserController) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSession, err := ctrl.session.CurrentSession.Get(r, ctrl.session.SessionName)
		if err != nil {
			until.WriteResponse(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		userSession.Options.MaxAge = -100
		if err := ctrl.session.CurrentSession.Save(r, w, userSession); err != nil {
			until.WriteResponse(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		until.WriteResponse(w, r, 200, map[string]interface{}{
			"result": "ok",
		}, nil)
	}
}
