package events

import (
	eventSerialize "ws/framework/plugin/event_serialize"
)

type TargetType int

const (
	TargetMedia TargetType = 1 // 点击+号
	TargetVoice TargetType = 3 // 按住录音按钮录音发这个。 因为录音发的是语音文件，所以发TargetMedia
	TargetText  TargetType = 8 // 点击输入框
)

// WamEventChatComposerAction 用户输入中状态
type WamEventChatComposerAction struct {
	WAMessageEvent

	ActionTarget float64
	ActionType   float64
}

type ChatComposerActionOption struct {
	Target TargetType
}

func WithChatComposerActionOption(target TargetType) ChatComposerActionOption {
	return ChatComposerActionOption{
		Target: target,
	}
}

func (event *WamEventChatComposerAction) InitFields(option interface{}) {
	event.ActionType = 1

	if opt, ok := option.(ChatComposerActionOption); ok {
		event.ActionTarget = float64(opt.Target)
	}
}

func (event *WamEventChatComposerAction) Serialize(buffer eventSerialize.IEventBuffer) {
	buffer.Header().
		SerializeNumber(event.Code, event.Weight)

	buffer.Body().
		SerializeNumber(0x1, event.ActionType)

	buffer.Footer().
		SerializeNumber(0x2, event.ActionTarget)
}
