package business

import (
	"github.com/google/uuid"
	waBinary "ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/application/constant/types"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	businessDB "ws/framework/application/data_storage/business/database"
)

// SMBQueryBusinessCollections .
type SMBQueryBusinessCollections struct {
	processor.BaseAction
	UserID string
	MagicV string
}

// RaiseErrorWhenNodeError 是否抛出node错误，如401，404等
func (m *SMBQueryBusinessCollections) RaiseErrorWhenNodeError() bool {
	return false
}

// Start .
func (m *SMBQueryBusinessCollections) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	sessionID := uuid.New().String()

	context.ResolveBusinessService().ContextExecute(func(profile *businessDB.BusinessProfile) {
		if len(profile.CatalogSessionID) > 0 {
			sessionID = profile.CatalogSessionID
		} else {
			profile.UpdateCatalogSessionID(sessionID)
		}
	})

	m.SendMessageId, err = context.SendNode(waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"id":      context.GenerateRequestID(),
			"smax_id": "35",
			"from":    context.ResolveJID(),
			"to":      types.ServerJID,
			"type":    message.IqGet,
			"xmlns":   "w:biz:catalog",
		},
		Content: []waBinary.Node{
			{
				Tag: "collections",
				Attrs: waBinary.Attrs{
					"biz_jid":          context.ResolveJID(),
					"empty_collection": "false",
				},
				Content: []waBinary.Node{
					{Tag: "after"},
					{Tag: "item_limit", Content: "3"},
					{Tag: "collection_limit", Content: "5"},
					{Tag: "width", Content: "144"},
					{Tag: "height", Content: "144"},
					{Tag: "catalog_session_id", Content: sessionID},
				},
			},
		},
	})
	return
}

// Receive .
func (m *SMBQueryBusinessCollections) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) error {
	next()

	return nil
}

// Error .
func (m *SMBQueryBusinessCollections) Error(context containerInterface.IMessageContext, err error) {
}
