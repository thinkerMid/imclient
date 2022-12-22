package scene

import (
	"unicode/utf8"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/action/common"
	"ws/framework/application/core/action/contact"
	privateChat "ws/framework/application/core/action/private_chat"
	privateChatCommon "ws/framework/application/core/action/private_chat/common"
	"ws/framework/application/core/action/user"
	"ws/framework/application/core/monitor"
	"ws/framework/application/core/processor"
	"ws/framework/lib/media_crypto"
	"ws/framework/plugin/media_decode/media_content"
)

// NewSession .
func NewSession(id string) Session {
	return Session{UserID: id}
}

// Session 私信场景
type Session struct {
	UserID      string
	ActionList  []containerInterface.IAction
	targetIsJID bool // 如果true 使用JID聊天基本是陌生人消息
}

// SetTargetIsJID .
func (c *Session) SetTargetIsJID(targetIsJID bool) {
	c.targetIsJID = targetIsJID
}

// Build .
func (c *Session) Build() containerInterface.IProcessor {
	if c.targetIsJID {
		return c.strangerSession()
	}

	return c.contactSession()
}

// 陌生人消息 不加会联系人 直接发送消息
func (c *Session) strangerSession() containerInterface.IProcessor {
	actionList := make([][]containerInterface.IAction, 0)

	// 仅限于没有陌生人记录的时候才会执行的
	actionList = append(actionList, []containerInterface.IAction{
		&contact.CheckStranger{UserID: c.UserID, FindByJID: c.targetIsJID, IsExist: false}, // 这里异常会跳出
		&contact.QueryStranger{UserID: c.UserID},
		&common.QueryAvatarUrl{UserID: c.UserID}, // 真机查的不是URL
		&common.QueryDisappearingMode{UserID: c.UserID},
		&common.SubscribeStatus{UserID: c.UserID},
		&common.QueryMultiDevicesIdentity{UserID: c.UserID}, // 这里一定是只获取设备0的会话信息
		&common.QueryUserDeviceList{UserID: c.UserID},       // 对方的设备列表
	})

	sendMsgActions := []containerInterface.IAction{
		&contact.CheckStranger{UserID: c.UserID, FindByJID: c.targetIsJID, IsExist: true}, // 这里异常会跳出
	}

	{
		// 如果只有一个action 并且不是发送消息的
		var unSendMessage bool
		if len(c.ActionList) == 1 {
			// 是不是输入状态
			_, unSendMessage = c.ActionList[0].(*common.InputChatState)
			if !unSendMessage {
				// 是不是标记已读
				_, unSendMessage = c.ActionList[0].(*privateChatCommon.ReceiveMessageMarkRead)
			}
		}

		// 不发送消息
		if unSendMessage {
			sendMsgActions = append(sendMsgActions, c.ActionList...)
		} else {
			// 要发消息
			sendMsgActions = append(sendMsgActions, &common.SubscribeStatus{UserID: c.UserID})
			sendMsgActions = append(sendMsgActions, &common.QueryMultiDevicesIdentity{UserID: c.UserID})
			sendMsgActions = append(sendMsgActions, &privateChatCommon.TrustedContactToken{UserID: c.UserID})
			sendMsgActions = append(sendMsgActions, c.ActionList...)
			sendMsgActions = append(sendMsgActions, &privateChatCommon.ChatWith{UserID: c.UserID})
		}
	}

	actionList = append(actionList, sendMsgActions)

	return processor.NewOnceComposeProcessor(
		actionList,
		processor.AliasName("strangerSession"),
		processor.AttachMonitor(&monitor.Session{}),
		processor.AttachMonitor(&monitor.PrivateChatMonitor{}),
		processor.AttachMonitor(&monitor.ReplaceUserIDWhenSession{DstUserID: c.UserID}),
	)
}

// 联系人消息 如果没有联系人会自动加
func (c *Session) contactSession() containerInterface.IProcessor {
	actionList := make([][]containerInterface.IAction, 0)

	// 没有添加过联系人
	actionList = append(actionList, []containerInterface.IAction{
		&contact.Check{UserID: c.UserID, FindByJID: c.targetIsJID, InAddressBook: false}, // 这里异常会跳出
		&contact.Query{UserID: c.UserID},
		&user.QueryStatusPrivacyList{IgnoreResponse: true},
		&contact.Add{UserID: c.UserID},
		&common.QueryAvatarUrl{UserID: c.UserID},
		&common.QueryDevicesIdentity{UserID: c.UserID, CheckHaveMultiDevice: true},
		&common.QueryUserDeviceList{UserID: c.UserID, CheckHaveMultiDevice: true},
	})

	actionList = append(actionList, []containerInterface.IAction{
		&common.SubscribeStatus{UserID: c.UserID},
	})

	// 没有会话
	actionList = append(actionList, []containerInterface.IAction{
		&contact.Check{UserID: c.UserID, FindByJID: c.targetIsJID, InAddressBook: true}, // 这里异常会跳出
		&common.QueryMultiDevicesIdentity{UserID: c.UserID},
	})

	sendMsgActions := []containerInterface.IAction{
		&contact.Check{UserID: c.UserID, InAddressBook: true, FindByJID: c.targetIsJID}, // 这里异常会跳出
		&privateChatCommon.TrustedContactToken{UserID: c.UserID},
	}
	sendMsgActions = append(sendMsgActions, c.ActionList...)
	sendMsgActions = append(sendMsgActions, &privateChatCommon.ChatWith{UserID: c.UserID})

	actionList = append(actionList, sendMsgActions)

	return processor.NewOnceComposeProcessor(
		actionList,
		processor.AliasName("contactSession"),
		processor.AttachMonitor(&monitor.Session{}),
		processor.AttachMonitor(&monitor.PrivateChatMonitor{}),
		processor.AttachMonitor(&monitor.ReplaceUserIDWhenSession{DstUserID: c.UserID}),
	)
}

// InputChatState .
func (c *Session) InputChatState(input bool) {
	c.ActionList = append(c.ActionList, &common.InputChatState{
		UserID: c.UserID,
		Input:  input,
	})
}

// MakeTextMessage .
func (c *Session) MakeTextMessage(messageText string) {
	// 非陌生人模式 需要加入模拟输入
	if !c.targetIsJID {
		inputLatency := utf8.RuneCountInString(messageText)

		c.ActionList = append(c.ActionList, &privateChatCommon.SimulateInputChatState{UserID: c.UserID, TotalInputLatency: uint32(inputLatency)})
	}

	c.ActionList = append(c.ActionList, &privateChat.SendText{
		UserID:      c.UserID,
		MessageText: messageText,
	})
}

// MakeTempMessage .
func (c *Session) MakeTempMessage(title, messageText, footer string, messageButton []privateChat.MessageButton) {
	c.ActionList = append(c.ActionList, &privateChat.SendTemp{
		UserID:      c.UserID,
		Title:       title,
		MessageText: messageText,
		Footer:      footer,
		Button:      messageButton,
	})
}

// MakeImageMessage .
func (c *Session) MakeImageMessage(content []byte, description string) error {
	imageContent, err := mediaContent.NewImageContent(content)

	if err != nil {
		return err
	}

	parser, ok := mediaCrypto.ParseMediaImage(content)
	if ok != nil {
		return ok
	}

	c.ActionList = append(c.ActionList, &privateChatCommon.SimulateInputChatState{UserID: c.UserID, TotalInputLatency: 1})
	c.ActionList = append(c.ActionList, &privateChat.SendImage{
		UserID:  c.UserID,
		Image:   imageContent,
		Parser:  *parser,
		Caption: description,
	})

	return nil
}

// MakeAudioMessage .
func (c *Session) MakeAudioMessage(content []byte) error {
	audioContent, err := mediaContent.NewAudioContent(content)

	if err != nil {
		return err
	}

	parser, err := mediaCrypto.ParseMediaVoice(content)
	if err != nil {
		return err
	}

	c.ActionList = append(c.ActionList, &privateChat.SendAudio{
		UserID: c.UserID,
		Audio:  audioContent,
		Parser: *parser,
	})

	return nil
}

// MakeVideoMessage .
func (c *Session) MakeVideoMessage(content []byte) error {
	videoContent, err := mediaContent.NewVideoContent(content)

	if err != nil {
		return err
	}

	parser, err := mediaCrypto.ParseMediaVideo(content)
	if err != nil {
		return err
	}

	c.ActionList = append(c.ActionList, &privateChatCommon.SimulateInputChatState{UserID: c.UserID, TotalInputLatency: 1})
	c.ActionList = append(c.ActionList, &privateChat.SendVideo{
		UserID: c.UserID,
		Video:  videoContent,
		Parser: *parser,
	})

	return nil
}

// MakeVCardMessage .
func (c *Session) MakeVCardMessage(contacts []string) {
	c.ActionList = append(c.ActionList, &privateChatCommon.VCardPhoneNumberCheck{Contacts: contacts})
	c.ActionList = append(c.ActionList, &privateChatCommon.SimulateInputChatState{UserID: c.UserID, TotalInputLatency: 1})
	c.ActionList = append(c.ActionList, &privateChat.SendVCard{
		UserID:   c.UserID,
		Contacts: contacts,
	})
}
