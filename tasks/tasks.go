package tasks

import (
	"fmt"
	"time"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/ws"
	"github.com/suvrick/go-kiss-server/services"
)

type TaskManager struct {
	userService *services.UserService
	botService  *services.BotService
	delay       uint64
}

func NewTaskManager(delay uint64, us *services.UserService, bs *services.BotService) *TaskManager {
	return &TaskManager{
		userService: us,
		botService:  bs,
		delay:       delay,
	}
}

func (t *TaskManager) Run() {
	fmt.Println(">>>>>>>>> Init task manager")

	for {
		time.Sleep(time.Minute * time.Duration(t.delay))

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
				t.botService.UpdateBot(*b)
			}
		}
	}

}

// Do запускает горутины
func (t *TaskManager) Do(bot *models.Bot) {
	gs := ws.NewSocket(bot)
	gs.Go()
}
