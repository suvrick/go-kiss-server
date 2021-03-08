package server

import (
	"github.com/suvrick/go-kiss-server/controllers"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Start ...
func Start(config *Config) error {

	db, err := createDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	router := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	proxyRepo := repositories.NewProxyRepository(db)
	botRepo := repositories.NewBotRepository(db)

	//middlewares.NewAuthMiddleWare(sg, userRepo)
	userService := services.NewUserService(userRepo)
	proxyService := services.NewProxyService(proxyRepo)
	botService := services.NewBotService(botRepo, userService, proxyService)

	controllers.NewUserController(router, userService)
	controllers.NewBotController(router, botService, userService)
	controllers.NewProxyController(router, proxyService)

	controllers.NewAdminController(router, userService)
	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "www/index.html")
	// })

	//return http.ListenAndServeTLS(":443", "../certs/cert.crt", "../certs/pk.key", srv)
	return router.Run(config.BindAddr)
}

func createDB(dbURL string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	// Migrate the schema
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Bot{})
	db.AutoMigrate(&model.Proxy{})

	return db, err
}
