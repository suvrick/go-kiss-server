package parser

import (
	"log"
	"testing"
)

var urls = []struct {
	SocialName string
	FrameUrl   string
}{
	{"mm", "https://bottle2.itsrealgames.com/mobile/build/v1593/?social_api=mm&type=mm&record_first_session=1&6=&is_app_user=1&session_key=95ac81be3e64f73fd56e4b12e004c9d0&vid=9914552860808032409&oid=9914552860808032409&app_id=543574&authentication_key=27861583a5de8bd7c1aa5955dfe757b8&session_expire=1623244207&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=2add120d8f235755ee8b34801798f775&window_id=CometName_659cc22f549d0d25dc769cb39daa105f&referer_type=left_menu&version=1593"},
	// {"vk", ""},
	// {"nn", "2312weqeqwe"},
	// {"vk", "https://bottle2.itsrealgames.com/www/vk.html?social_api=vk&type=vk&record_first_session=1&6&api_url=https://api.vk.com/api.php&api_id=1930071&api_settings=8463&viewer_id=579349535&viewer_type=2&sid=7bd6b8b4252d4c39f1d8506bc364c9d9beea5e333121310f7e198f05bf513228629ffc2a770dbc0576f9e&secret=fa5abbef74&access_token=84d52e17d6032b1c9c9463720d72afec64dd98a94a5bf11310be158cd715b13296a5656b83529245e9a63&user_id=579349535&group_id=0&is_app_user=1&auth_key=39c574533868fe0a5e6ca2ba39f49495&language=0&parent_language=0&is_secure=1&stats_hash=b3e71399bfbb87a7f4&ads_app_id=1930071_4bb280b66d445bceae&referrer=left_nav&lc_name=3d8602a2&platform=web&whitelist_scopes=friends,photos,video,stories,pages,status,notes,wall,docs,groups,stats,market,ads,notifications&group_whitelist_scopes=stories,photos,app_widget,messages,wall,docs,manage&is_widescreen=0&hash=#api=vk&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/&width=1000&height=690&locale_version=62&useApiType=vk&"},
	// {"sa", "https://bottle2.itsrealgames.com/www/sa.html?time=1615887387318&&userId=120326146&sessionKey=caa2ca27d6f5c52c737c587c2800df36&authKey=e03daaa0da0cb6ec09b715b5287d098c&net_type=32&useApiType=sa&email=c9999cc@mail.ru&locale=RU&#time=1615887387318&userId=120326146&sessionKey=caa2ca27d6f5c52c737c587c2800df36&authKey=e03daaa0da0cb6ec09b715b5287d098c&net_type=32&useApiType=sa&email=c9999cc@mail.ru&locale=RU&api=sa&packageName=bottlePackage&config=config_release.xml&protocol=https:&international=false&locale_url=../resources/locale/&width=1000&height=690&locale_version=62&"},
	// {"ok", "https://bottle2.itsrealgames.com/www/ok.html?api=ok&6&record_first_session=1&container=true&web_server=https%3A%2F%2Fok.ru&first_start=0&logged_user_id=910908306174&sig=76e32e621230b230715484feeba46060&refplace=vitrine_user_apps_portlet&new_sig=1&apiconnection=83735040_1615887597845&authorized=1&session_key=-s-d274zznZ0fc7vwsb1d02Lwrc0d5cPyS5049eO-x3z3b9vzQZca7cS0r102-aP.q1z242y0N1a71fxxP71703SxS3Yb20Ouw1b1b3Ovqa&clientLog=0&session_secret_key=08edc349e34a428481632995e44740d0&auth_sig=c7ea9dd8af90d7cbceb223674323b48d&api_server=https%3A%2F%2Fapi.ok.ru%2F&ip_geo_location=RU%2C91%2CAchinsk&application_key=CBADLOPFABABABABA#api=ok&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/&width=1000&height=690&locale_version=62&useApiType=ok&"},
	// {"mm", "https://bottle2.itsrealgames.com/www/mm.html?is_app_user=1&session_key=cf5910cb442740a0169d8489281296ad&vid=6656435751632694779&oid=6656435751632694779&app_id=543574&authentication_key=2fb5074823068fea24d3287bbc92c8ff&session_expire=1615974154&ext_perm=notifications%2Cemails%2Cpayments&sig=a7d04fa68a2f55c8334d0126607212b7&window_id=CometName_d5afdf7e6d4b27694325ce9aa5c33256&referer_type=left_menu#api=mm&packageName=bottlePackage&config=config_release.xml&protocol=https:&locale=RU&international=false&locale_url=../resources/locale/&width=1000&height=690&locale_version=62&useApiType=mm&"},
	// {"mm", "https://bottle2.itsrealgames.com/mobile/build/v1476/?is_app_user=1&session_key=0964fca2d28794cfd478f79d63dbae0e&vid=6986266616779218243&oid=6986266616779218243&app_id=543574&authentication_key=71c9c773a338b0eee5ccc1debbcea74e&session_expire=1615958987&ext_perm=photos%2Cfriends%2Cevents%2Cguestbook%2Cmessages%2Cnotifications%2Cstream%2Cemails%2Cpayments&sig=65046c3169d1701fd73745a41f54a772&window_id=CometName_cd8ba653f2412736dc03696ccfe76b1e&referer_type=left_menu&prev_url=https%3A%2F%2Fbottle2.itsrealgames.com%2Fwww%2Fmm.html%3Fis_app_user%3D1%26session_key%3D0964fca2d28794cfd478f79d63dbae0e%26vid%3D6986266616779218243%26oid%3D6986266616779218243%26app_id%3D543574%26authentication_key%3D71c9c773a338b0eee5ccc1debbcea74e%26session_expire%3D1615958987%26ext_perm%3Dphotos%252Cfriends%252Cevents%252Cguestbook%252Cmessages%252Cnotifications%252Cstream%252Cemails%252Cpayments%26sig%3D65046c3169d1701fd73745a41f54a772%26window_id%3DCometName_cd8ba653f2412736dc03696ccfe76b1e%26referer_type%3Dleft_menu&version=1476"},
}

func TestInitialize(t *testing.T) {
	Initialize()
	if len(parseParams) == 0 {
		t.Errorf("Error load parser/config.json")
	}
}

func TestParse(t *testing.T) {

	for _, u := range urls {
		loginData := NewLoginParams(u.FrameUrl)
		log.Println(loginData.ToString())

		if loginData.SocialName != u.SocialName {
			t.Errorf(loginData.ToString())
		}
	}
}
