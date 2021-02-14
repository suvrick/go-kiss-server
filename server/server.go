package server

import (
	"database/sql"
	"net/http"

	"github.com/suvrick/go-kiss-server/controllers"
	"github.com/suvrick/go-kiss-server/middlewares"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/store"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq" // ...
)

// Start ...
func Start(config *Config) error {

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	router := mux.NewRouter()

	sg := session.NewSessionGame(config.SessionKey)

	userRepo := store.NewUserRepository(db)
	proxyRepo := store.NewProxyRepository(db)

	middlewares.NewAuthMiddleWare(sg, userRepo)

	controllers.NewUserController(router, userRepo, sg)
	controllers.NewProxyController(router, proxyRepo)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/index.html")
	})

	//return http.ListenAndServeTLS(":443", "../certs/cert.crt", "../certs/pk.key", srv)
	return http.ListenAndServe(config.BindAddr, router)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
