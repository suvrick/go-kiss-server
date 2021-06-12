package tasks

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/ws"
	"github.com/suvrick/go-kiss-server/repositories"
	"github.com/suvrick/go-kiss-server/services"
)

type TaskManager struct {
	userService  *services.UserService
	botService   *services.BotService
	proxyManager *repositories.ProxyRepository
	delay        uint64
}

func NewTaskManager(delay uint64, us *services.UserService, bs *services.BotService, pr *repositories.ProxyRepository) *TaskManager {
	return &TaskManager{
		userService:  us,
		botService:   bs,
		proxyManager: pr,
		delay:        delay,
	}
}

func (t *TaskManager) Run() {

	fmt.Println(">>>>>>>>> Init task manager")
	skip := false

	for {

		fmt.Printf("Now: %v\n", time.Now().Format("02.01.2006 15:04:05"))

		if time.Now().Hour() != 1 || skip {
			skip = false
			d := time.Minute * time.Duration(t.delay)
			fmt.Printf("Sleep delay %v\n", d)
			time.Sleep(d)
			continue
		}

		skip = true
		fmt.Println("Reset  proxy!")
		t.proxyManager.ResetProxy()

		fmt.Println(">>>>>>>>> Start loop tasks")

		users, err := t.userService.AllUser()

		if err != nil {
			return
		}

		for _, u := range users {
			fmt.Printf(">>>>>>>>> Get userID %v\n", u.ID)
			bots, err := t.botService.AllByUserID(u.ID)

			if err != nil {
				continue
			}

			fmt.Printf(">>>>>>>>> Get bots %v\n", len(bots))
			for _, b := range bots {
				t.Do(b)
			}
		}
	}

}

// Do запускает горутины
func (t *TaskManager) Do(bot *models.Bot) {
	gs := ws.NewSocket(bot)
	gs.SetProxyManager(t.proxyManager)
	gs.Go()
	t.botService.UpdateBot(bot)
}
