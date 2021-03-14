package server

import (
	"github.com/suvrick/go-kiss-server/controllers"
	"github.com/suvrick/go-kiss-server/middlewares"
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
	router.Use(middlewares.CORSMiddleware())

	userRepo := repositories.NewUserRepository(db)
	proxyRepo := repositories.NewProxyRepository(db)
	botRepo := repositories.NewBotRepository(db)
	kissRepo := repositories.NewAutoKissRepository(db)
	stateRepo := repositories.NewStateDowloadRepository(db)

	//middlewares.NewAuthMiddleWare(sg, userRepo)
	userService := services.NewUserService(userRepo)
	proxyService := services.NewProxyService(proxyRepo)
	botService := services.NewBotService(botRepo, userService, proxyService)
	kissService := services.NewAutoKissService(kissRepo)
	stateDownloadService := services.NewStateDownloadService(stateRepo)

	controllers.NewAutoKissController(router, kissService, stateDownloadService)

	controllers.NewUserController(router, userService)
	controllers.NewBotController(router, botService, userService)
	controllers.NewProxyController(router, proxyService)

	controllers.NewAdminController(router, userService, kissService, stateDownloadService)

	router.GET("", func(c *gin.Context) {
		c.File("www/autokiss/index.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("www/login.html")
	})

	router.GET("/register", func(c *gin.Context) {
		c.File("www/register.html")
	})

	router.GET("/admin", func(c *gin.Context) {
		c.File("www/admin.html")
	})

	//return http.ListenAndServeTLS(":443", "../certs/cert.crt", "../certs/pk.key", router)
	return router.Run(config.BindAddr)
}

func createDB(dbURL string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	// Migrate the schema
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Bot{})
	db.AutoMigrate(&model.Proxy{})
	db.AutoMigrate(&model.KissUser{})
	db.AutoMigrate(&model.KissState{})

	return db, err
}
