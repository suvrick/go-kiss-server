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
}

func NewTaskManager(us *services.UserService, bs *services.BotService) *TaskManager {
	return &TaskManager{
		userService: us,
		botService:  bs,
	}
}

func (t *TaskManager) Run() {

	for {

		time.Sleep(time.Minute * 5)

		fmt.Println("Start task manager")

		users, err := t.userService.AllUser()

		if err != nil {
			return
		}

		for _, u := range users {
			fmt.Printf("Get userID %v\n", u.ID)
			bots, err := t.botService.AllByUserID(u.ID)

			if err != nil {
				continue
			}

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
