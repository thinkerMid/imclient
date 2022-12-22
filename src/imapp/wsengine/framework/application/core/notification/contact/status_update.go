package contactNotification

import (
	"ws/framework/application/constant"
	containerInterface "ws/framework/application/container/abstract_interface"
	messageResultType "ws/framework/application/core/result/constant"
	"ws/framework/external"
	functionTools "ws/framework/utils/function_tools"
)

// UnknownCharacter 不知道为什么快捷设置的状态前3个字节是奇怪的字符
var UnknownCharacter = []byte{226, 128, 142}

// SignatureUpdate 签名更新
type SignatureUpdate struct{}

// Receive .
func (s SignatureUpdate) Receive(context containerInterface.IMessageContext) (err error) {
	node := context.Message()

	ag := node.AttrGetter()
	if ag.String("type") != "status" {
		return
	}

	child, ok := node.GetOptionalChildByTag("set")
	if !ok {
		return
	}

	childAttrGetter := child.AttrGetter()

	/**
	[S] <notification from="85256048573@s.whatsapp.net" id="776126216" notify="hhhhggffd" offline="1" t="1661304251" type="status"><set>e2808ee59ca8e79c8be794b5e5bdb1</set></notification>
	[S] <notification from="85296475450@s.whatsapp.net" id="3925341735" notify="hhhhggffd" offline="1" t="1661304251" type="status"><set hash="Fg8T"/></notification>
	*/
	if len(childAttrGetter.String("hash")) == 0 {
		jid := ag.JID("from")
		var signatureText string

		bStatus, ok := child.Content.([]byte)
		if ok {
			if len(bStatus) > len(UnknownCharacter) && functionTools.SliceEqual(bStatus[:len(UnknownCharacter)], UnknownCharacter) {
				bStatus = bStatus[len(UnknownCharacter):]
			}

			signatureText = string(bStatus)
		}

		context.AppendResult(containerInterface.MessageResult{
			ResultType: messageResultType.ContactSignatureUpdate,
			IContent: external.ProfileUpdate{
				JIDNumber: jid.User,
				Content:   signatureText,
			},
		})
	}

	return constant.AbortedError
}
