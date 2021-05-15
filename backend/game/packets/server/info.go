package server

import (
	"io"

	"github.com/suvrick/go-kiss-server/game/packets/decode"
)

/*
	public nid: string; // PlayerInfoParser.NET_ID; [I;
	public type: NetType; // PlayerInfoParser.TYPE; [B;
	public name: string; // PlayerInfoParser.NAME; [S;
	public sex: Gender; // PlayerInfoParser.SEX; [B;
	public tag: number; // PlayerInfoParser.TAG; [I;
	public referrer: number; // PlayerInfoParser.REFERRER; [I;
	public bdate: number; // PlayerInfoParser.BDAY; [I;
	public avatar: string; public avatar_status: number; // PlayerInfoParser.PHOTO; [SB;
	public profile: string; // PlayerInfoParser.PROFILE; [S;

	public status: string; // PlayerInfoParser.STATUS; [S;
	public countryId: number; // PlayerInfoParser.COUNTRY_ID; [B;
	public online: boolean; // PlayerInfoParser.ONLINE; [B;
	public admirer_id: number; // PlayerInfoParser.ADMIRER_ID; [I;
	public admirer_price: number; // PlayerInfoParser.ADMIRER_PRICE; [I;
	public admirer_time_finish: number; // PlayerInfoParser.ADMIRER_LEFT; [I; it is timestamp
	public views: number; // PlayerInfoParser.VIEWS; [I;
	public vip: boolean; // PlayerInfoParser.IS_VIP; [B;
	public color: number; // PlayerInfoParser.COLOR; [B;
	public kisses: number; // PlayerInfoParser.KISSES; [I;
	public hearts: number; // PlayerInfoParser.HEARTS; [I;
	public gifts: number; // PlayerInfoParser.GIFTS; [I;
	public lastGifts: [number, number, number, string][]; // PlayerInfoParser.PLAYER_GIFTS; ["[IIBS]"; // [source_id:I; gift_id:I; is_private:B; message:S]
	public device: DeviceType; // PlayerInfoParser.DEVICE; [B;
	public wedding_id: number; // PlayerInfoParser.WEDDING_ID; [I;
	public achievements: [number, number, number][]; // PlayerInfoParser.ACHIEVEMENTS; ["[III]";
	public room_id: number; // PlayerInfoParser.ROOM_ID; [I;
	public collections: [number, number][]; // PlayerInfoParser.COLLECTIONS_SETS; ["[BI]";
	public avatar_id: number; // PlayerInfoParser.AVATAR_ID; [B;
	public rights: number; // PlayerInfoParser.RIGHTS; [B;
	public register_time: number; // PlayerInfoParser.REGISTER_TIME; [I;
	public logout_time: number; // PlayerInfoParser.LOGOUT_TIME; [I;
	public photos: string[]; public photos_statuses: number[]; // PlayerInfoParser.PHOTOS; ["[S][B]";
*/

type InfoServerPacket struct {
	Field1        int32
	Field2        int16
	Field3        int32
	NetId         int64
	NetType       byte
	Name          string
	Sex           byte
	Tag           int32
	Referrer      int32
	Ddate         int32
	Avatar        string
	Avatar_status byte
	Profile       string
	Status        string
	CountryId     byte
	Online        byte
	Blot_time     int32
	Admirer_id    int32
	Admirer_price int32
	Admirer_time  int32
	Views         int32
	Vip           byte
	// Color         byte
	// Kisses        int32
	// Hearts        int32
	// Gifts         int32
}

// Parse ...
func (pack *InfoServerPacket) Parse(buffer io.Reader) error {
	return decode.Fill(pack, buffer)
}
