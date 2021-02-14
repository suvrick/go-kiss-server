package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/store"
	"github.com/suvrick/go-kiss-server/until"

	"github.com/gorilla/mux"
)

// ProxyController ...
type ProxyController struct {
	router          *mux.Router
	proxyRepository *store.ProxyRepository
}

// NewProxyController ...
func NewProxyController(router *mux.Router, proxyRepository *store.ProxyRepository) *ProxyController {

	ctrl := &ProxyController{
		router:          router,
		proxyRepository: proxyRepository,
	}

	proxy := ctrl.router.PathPrefix("/proxy").Subrouter()

	proxy.Use(middlewares.AuthMiddlewareInstance.Do)

	proxy.HandleFunc("/add", ctrl.addProxy()).Methods("POST")
	proxy.HandleFunc("/all", ctrl.allProxy()).Methods("GET")
	proxy.HandleFunc("/find", ctrl.findProxy()).Methods("GET")
	proxy.HandleFunc("/free", ctrl.freeProxy()).Methods("GET")
	return ctrl
}

// POST add proxy
func (pc *ProxyController) addProxy() http.HandlerFunc {

	type request struct {
		Proxies []string `json: proxies`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := until.JSONBind(r, req); err != nil {
			until.WriteResponse(w, r, http.StatusBadRequest, nil, err)
			return
		}

		proxies := req.Proxies
		result := make([]int, 0)
		for _, item := range proxies {
			p := model.NewProxy(item)

			id, err := pc.proxyRepository.Create(p)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			result = append(result, id)
		}

		until.WriteResponse(w, r, 200, result, nil)
	}
}

// GET all proxy
func (pc *ProxyController) allProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := pc.proxyRepository.All()
		until.WriteResponse(w, r, http.StatusOK, p, err)
	}
}

// GET find proxy
func (pc *ProxyController) findProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		val, ok := r.URL.Query()["id"]

		if !ok {
			until.WriteResponse(w, r, http.StatusOK, nil, errors.New("Bad param"))
			return
		}

		id, _ := strconv.ParseInt(val[0], 10, 32)

		p, err := pc.proxyRepository.Find(int(id))
		if err != nil {
			until.WriteResponse(w, r, http.StatusOK, nil, err)
			return
		}

		until.WriteResponse(w, r, http.StatusOK, p, nil)
	}
}

// GET free proxy
func (pc *ProxyController) freeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := pc.proxyRepository.Free()
		until.WriteResponse(w, r, http.StatusOK, p, err)
	}
}
