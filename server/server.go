package server

import (
	"log"
	"net/http"

	"github.com/suvrick/go-kiss-server/controllers"
	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/model"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/services"
	"github.com/suvrick/go-kiss-server/session"
	"github.com/suvrick/go-kiss-server/tasks"

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

	setStaticFile(router)

	userRepo := repositories.NewUserRepository(db)
	botRepo := repositories.NewBotRepository(db)

	session.SetDb(userRepo)

	userService := services.NewUserService(userRepo)
	botService := services.NewBotService(botRepo, userService)

	controllers.NewUserController(router, userService)
	controllers.NewBotController(router, botService, userService)
	controllers.NewAdminController(router, userService)

	taskServer := tasks.NewTaskManager(60*6, userService, botService)
	go taskServer.Run()

	return http.ListenAndServeTLS(":443", "../certs/cert.crt", "../certs/pk.key", router)
	//return router.Run(config.BindAddr)
}

func setStaticFile(router *gin.Engine) {
	router.StaticFile("/", "./www/index.html")
	router.StaticFile("/login", "./www/login.html")

	router.StaticFile("autokiss.zip", "./www/autokiss/autokiss.zip")
	router.Static("bootstrap", "./www/bootstrap")
	router.Static("css", "./www/css")
	router.Static("js", "./www/js")
}

func createDB(dbURL string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if !db.Migrator().HasTable("bots") {
		if err := db.Migrator().CreateTable(&models.Bot{}); err != nil {
			log.Fatalln(err.Error())
		}
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&models.Bot{})
	db.AutoMigrate(&models.LoggerLine{})

	return db, err
}
