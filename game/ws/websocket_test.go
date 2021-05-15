package ws_test

import (
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/ws"
)

func TestNewSocket(t *testing.T) {
	u := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000935&sessionKey=e0a663c461473539f07c3dadc486543&authKey=9d39455de0b34d1f2dfcf4390523cf19&net_type=32&useApiType=sa"
	b := models.NewBot(u)
	s := ws.NewSocket(b)
	s.Go()

	fmt.Println(b.ToString())
	b.PrintLog()
}
