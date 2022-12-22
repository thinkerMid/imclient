package monitor

import (
	containerInterface "ws/framework/application/container/abstract_interface"
	privateChat "ws/framework/application/core/action/private_chat"
	privateChatCommon "ws/framework/application/core/action/private_chat/common"
	. "ws/framework/application/core/wam"
	"ws/framework/application/core/wam/events"
)

// PrivateChatMonitor .
type PrivateChatMonitor struct {
	sendStateCount int32
}

// OnStart .
func (p *PrivateChatMonitor) OnStart(ioc containerInterface.IAppIocContainer) {
	LogManager().SwitchAppMenu(ioc, PageSession)
}

// OnActionStartBefore .
func (p *PrivateChatMonitor) OnActionStartBefore(_ interface{}, _ containerInterface.IMessageContext) {
}

// OnActionStartAfter .
func (p *PrivateChatMonitor) OnActionStartAfter(_ interface{}, _ containerInterface.IMessageContext) {
}

func (p *PrivateChatMonitor) OnActionStartFail(action interface{}, context containerInterface.IMessageContext) {

}

func (p *PrivateChatMonitor) OnActionStartSuccess(action interface{}, context containerInterface.IMessageContext) {
}

// ActionExecuteSuccess .
func (p *PrivateChatMonitor) ActionExecuteSuccess(action interface{}, context containerInterface.IMessageContext) {
	switch action.(type) {
	case *privateChatCommon.SimulateInputChatState:
		state := action.(*privateChatCommon.SimulateInputChatState)
		p.sendStateCount = state.TotalCount
	case *privateChat.SendText:
		text := action.(*privateChat.SendText)
		contact := context.ResolveContactService().FindByJID(text.UserID)
		if contact != nil {
			LogManager().LogSendText(context, p.sendStateCount, !contact.ChatWith)
		}
	case *privateChat.SendImage:
		req := action.(*privateChat.SendImage)
		img := req.Image
		parser := req.Parser

		media := Media{
			Image: &Image{
				Width:     int32(img.Width),
				Height:    int32(img.Height),
				Size:      int32(parser.File.FileLength),
				FirstScan: parser.FirstScanLength,
				LowScan:   parser.LowQualityScanLength,
				MidScan:   parser.MidQualityScanLength,
			},
		}
		LogManager().LogSendMedia(context, events.MediaImage, media)
	case *privateChat.SendVideo:
		req := action.(*privateChat.SendVideo)
		vdo := req.Video
		parser := req.Parser

		media := Media{
			Video: &Video{
				Width:  int32(vdo.Width),
				Height: int32(vdo.Height),
				Size:   int32(parser.FileLength),
			},
		}
		LogManager().LogSendMedia(context, events.MediaVideo, media)
	case *privateChat.SendAudio:
		req := action.(*privateChat.SendAudio)
		parser := req.Parser

		media := Media{
			Voice: &Voice{
				Size: int32(parser.FileLength),
			},
		}
		LogManager().LogSendMedia(context, events.MediaVoice, media)

	case *privateChat.SendVCard:
		req := action.(*privateChat.SendVCard)
		LogManager().LogSendVCards(context, req.Contacts)
	}
}

func (p *PrivateChatMonitor) ActionExecuteFailure(action interface{}, context containerInterface.IMessageContext) {

}

// OnExit .
func (p *PrivateChatMonitor) OnExit(containerInterface.IAppIocContainer) {

}
