package ws_test

import (
	"fmt"
	"testing"

	"github.com/suvrick/go-kiss-server/game/models"
	"github.com/suvrick/go-kiss-server/game/ws"
)

func TestNewSocket(t *testing.T) {
	//u := "https://bottle2.itsrealgames.com/www/sa.html?time=1561991107710&&userId=1000935&sessionKey=e0a663c461473539f07c3dadc486543&authKey=9d39455de0b34d1f2dfcf4390523cf19&net_type=32&useApiType=sa"
	//u := "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=8207&viewer_id=196548541&viewer_type=2&sid=8b72392f19d79c0770b22ff12713d162cf1f32730acb0a10571b908471a2ce692101aacf0a8fcefd3945c&secret=24e9ee7629&access_token=958bd325c1fd16196be46d33be9a93f5d353233cf9afc51b0bf63a544cc11c1029d16fe2523f736eecdaa&user_id=196548541&group_id=0&is_app_user=1&auth_key=70691bedeac93e2ffa7fe65e3cbed943&language=0&parent_language=0&is_secure=1&stats_hash=3fca0f39cdf48e2369&ads_app_id=1930071_ea8c4ceaca3ace2676&referrer=menu&lc_name=276dd73c&platform=web&hash="
	u := "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=95ac81be3e64f73fd56e4b12e004c9d0&vid=9914552860808032409&oid=9914552860808032409&app_id=543574&authentication_key=27861583a5de8bd7c1aa5955dfe757b8&session_expire=1623244207&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=2add120d8f235755ee8b34801798f775&window_id=CometName_659cc22f549d0d25dc769cb39daa105f&referer_type=left_menu&version=1593"
	b := models.NewBot(u)
	s := ws.NewSocket(b)
	s.Go()

	fmt.Println(b.ToString())
	b.PrintLog()
}
