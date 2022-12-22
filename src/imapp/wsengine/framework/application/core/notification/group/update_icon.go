package groupNotification

import (
	"ws/framework/application/constant"
	"ws/framework/application/container/abstract_interface"
	groupDB "ws/framework/application/data_storage/group/database"
)

/**
<notification from="120363042017262679@g.us" id="3634898369" notify="uuu" t="1650612224" type="picture">
<set author="85256048573@s.whatsapp.net" id="1650612224" jid="120363042017262679@g.us"/>
</notification>

<notification from="120363042017262679@g.us" id="329286575" notify="uuu" t="1650612254" type="picture">
<delete author="85256048573@s.whatsapp.net" jid="120363042017262679@g.us"/>
</notification>
*/

// UpdateIcon .
type UpdateIcon struct{}

// Receive .
func (m UpdateIcon) Receive(context containerInterface.IMessageContext) error {
	groupId, err := parseGroupNotification(context.Message())
	if err != nil {
		return nil
	}

	if context.Message().AttrGetter().String("type") != "picture" {
		return nil
	}

	_, ok := context.Message().GetOptionalChildByTag("set")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateHaveGroupIcon(true)
		})

		return constant.AbortedError
	}

	_, ok = context.Message().GetOptionalChildByTag("delete")
	if ok {
		context.ResolveGroupService().ContextExecute(groupId, func(group *groupDB.Group) {
			group.UpdateHaveGroupIcon(false)
		})

		return constant.AbortedError
	}

	return nil
}
