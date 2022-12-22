package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type ContactAction int

const (
	ContactOpen ContactAction = iota + 1
	ContactClose
	ContactCancel
)

type OpenType int

const (
	OpenContactAdd OpenType = iota + 1
	OpenContactView
)

type WamEventIphoneAddContactEvent struct {
	WAMessageEvent

	ActionType     float64 // 1: 打开 2: 关闭添加联系人 3：取消
	DuplicateEvent float64
	EventType      float64 // 添加联系人进入:1 查看联系人进入:2
	ContactSource  float64 // 添加联系人的方式，写死5就好 TODO:看到个2
	Business       float64 // 非商务号，写死0
	SaveSuccess    float64
	SessionId      float64
}

type IphoneAddContactEventOption struct {
	Action   ContactAction
	Session  float64
	OpenType OpenType
}

func WithIphoneAddContactEvent(act ContactAction, session float64, open OpenType) IphoneAddContactEventOption {
	return IphoneAddContactEventOption{
		Action:   act,
		Session:  session,
		OpenType: open,
	}
}

func (event *WamEventIphoneAddContactEvent) InitFields(option interface{}) {
	if opt, ok := option.(IphoneAddContactEventOption); ok {
		event.ActionType = float64(opt.Action)
		event.SessionId = opt.Session
		event.EventType = float64(opt.OpenType)

		switch opt.OpenType {
		case OpenContactAdd:
			event.ContactSource = 2
		case OpenContactView:
			event.ContactSource = 5
		}
	}

	switch event.ActionType {
	case 1: // 打开
		event.DuplicateEvent = 0
		event.SaveSuccess = 0
	case 2: // 添加成功时关闭
		event.DuplicateEvent = 1
		event.SaveSuccess = 1
	case 3: // 添加失败时关闭
		event.DuplicateEvent = 0
		event.SaveSuccess = 0
	}

	event.Business = 0.000000
}

func (event *WamEventIphoneAddContactEvent) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x2, event.ActionType).
		SerializeNumber(0x7, event.DuplicateEvent).
		SerializeNumber(0x3, event.EventType).
		SerializeNumber(0x5, event.Business).
		SerializeNumber(0x4, event.SaveSuccess).
		SerializeNumber(0x6, event.SessionId)

	buffer.Footer().
		SerializeNumber(0x1, event.ContactSource)
}
