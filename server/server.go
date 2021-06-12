package server

import (
	"fmt"
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

	router.GET("/", func(c *gin.Context) {
		c.File("../www/index.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("../www/login.html")
	})

	router.Static("bootstrap", "../www/bootstrap")
	router.Static("css", "../www/css")
	router.Static("js", "../www/js")

	//setStaticFile(router)

	userRepo := repositories.NewUserRepository(db)
	botRepo := repositories.NewBotRepository(db)
	proxyRepo := repositories.NewProxyRepository(db)

	session.SetDb(userRepo)

	userService := services.NewUserService(userRepo)
	botService := services.NewBotService(botRepo, userService, proxyRepo)

	controllers.NewUserController(router, userService)
	controllers.NewAdminController(router, userService)
	controllers.NewWsController(router, userService, botService)

	router.POST("proxy/add", func(c *gin.Context) {
		type FormData struct {
			Proxies []string `json:"proxies"`
		}

		data := FormData{}

		if err := c.ShouldBindJSON(&data); err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, u := range data.Proxies {
			p, err := repositories.NewProxy(u)
			if err != nil {
				continue
			}

			proxyRepo.Add(p)
		}

	})

	taskServer := tasks.NewTaskManager(60, userService, botService, proxyRepo)
	go taskServer.Run()

	return http.ListenAndServeTLS(":443", "../certs/cert.crt", "../certs/pk.key", router)
	//return router.Run(config.BindAddr)
}

func setStaticFile(router *gin.Engine) {

	router.StaticFile("/", "./www/index.html")
	router.StaticFile("/login", "./www/login.html")

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

	db.AutoMigrate(&repositories.Proxy{})

	return db, err
}
