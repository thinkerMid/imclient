package events

import (
	"fmt"
	eventSerialize "ws/framework/plugin/event_serialize"
)

/*
* 账号注册登陆时构造
* 一次注册登陆构造四个包，@QueryType顺序2-0-2-1，对应回复@ResponseParams
 */
var (
	ResponseParams = []string{
		"cat=all&lg=%v-%v&ver=2&country=%v",
		"img=d2dc590fc2ffafb664a7541f409defff",
		"id=whatsappcuppy&lg=%v",
	}
)

type StickerCommonQueryToStaticServerOption struct {
	QueryType int
	Language  string
	Country   string
}

func WithStickerCommonQueryToStaticServerOption(query int, language, country string) StickerCommonQueryToStaticServerOption {
	return StickerCommonQueryToStaticServerOption{
		QueryType: query,
		Language:  language,
		Country:   country,
	}
}

type WamEventStickerCommonQueryToStaticServer struct {
	WAMessageEvent

	ResponseCode float64
	Params       string
	QueryType    float64
}

// InitFields .
func (event *WamEventStickerCommonQueryToStaticServer) InitFields(option interface{}) {
	if opt, ok := option.(StickerCommonQueryToStaticServerOption); ok {
		event.QueryType = float64(opt.QueryType)
		event.Params = ResponseParams[opt.QueryType]

		switch opt.QueryType {
		case 0:
			event.Params = fmt.Sprintf(event.Params, opt.Language, opt.Country, opt.Country)
		case 2:
			event.Params = fmt.Sprintf(event.Params, opt.Language)
		}
	}
	event.ResponseCode = 200
}

func (event *WamEventStickerCommonQueryToStaticServer) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x2, event.ResponseCode).
		SerializeString(0x3, event.Params)

	buffer.Footer().
		SerializeNumber(0x1, event.QueryType)
}
