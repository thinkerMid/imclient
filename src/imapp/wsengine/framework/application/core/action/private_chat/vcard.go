package privateChat

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"google.golang.org/protobuf/proto"
	"regexp"
	waProto "ws/framework/application/constant/binary/proto"
	"ws/framework/application/constant/types"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	"ws/framework/application/core/result/constant"
	"ws/framework/lib/msisdn"
	"ws/framework/utils"
	"ws/framework/utils/xmpp"
)

const (
	vCardTemplate string = "BEGIN:VCARD\n" +
		"VERSION:3.0\n" +
		"N:;;;;\n" +
		"FN:%v\n" +
		"TEL;type=CELL;type=VOICE;waid=%v:%v\n" +
		"END:VCARD"
)

// SendVCard .
type SendVCard struct {
	processor.BaseAction
	UserID   string
	Contacts []string
}

// Start .
func (c *SendVCard) Start(context containerInterface.IMessageContext, _ containerInterface.NextActionFn) (err error) {
	message := waProto.Message{
		MessageContextInfo: &waProto.MessageContextInfo{
			DeviceListMetadata: &waProto.DeviceListMetadata{
				SenderTimestamp: proto.Uint64(uint64(utils.GetCurTime())),
			},
			DeviceListMetadataVersion: proto.Int32(2),
		},
	}

	err = formatVCardMessageNode(&message, c.Contacts)
	if err != nil {
		return err
	}

	dstJID := types.NewJID(c.UserID, types.DefaultUserServer)
	node := xmpp.CreateMessageNode(dstJID, xmpp.MediaMessageType)

	mediaType := xmpp.ContactCard
	if len(c.Contacts) > 1 {
		mediaType = xmpp.ContactCardArray
	}

	err = encodeProtocolMessage(context, dstJID, &node, mediaType, &message)
	if err != nil {
		return err
	}

	c.SendMessageId, err = context.SendNode(node)

	return err
}

// Receive .
func (c *SendVCard) Receive(context containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatVCard,
		Content:    c.SendMessageId,
	})

	next()

	return
}

// Error .
func (c *SendVCard) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.SendPrivateChatVCard,
		Error:      err,
	})
}

func formatVCardMessageNode(msg *waProto.Message, contacts []string) (err error) {
	if len(contacts) == 1 {
		international, telephone, err := castPhoneNumber(contacts[0])
		if err != nil {
			return err
		}

		content := fmt.Sprintf(vCardTemplate, international, telephone, international)

		msg.ContactMessage = &waProto.ContactMessage{
			DisplayName: proto.String(international),
			Vcard:       proto.String(content),
		}
	} else if len(contacts) > 1 {
		var array []*waProto.ContactMessage

		for _, contact := range contacts {
			international, telephone, err := castPhoneNumber(contact)
			if err != nil {
				continue
			}

			content := fmt.Sprintf(vCardTemplate, international, telephone, international)

			array = append(array, &waProto.ContactMessage{
				DisplayName: proto.String(international),
				Vcard:       proto.String(content),
			})
		}

		msg.ContactsArrayMessage = &waProto.ContactsArrayMessage{
			DisplayName: proto.String(fmt.Sprintf("%v 位联系人", len(array))),
			Contacts:    array,
		}
	}
	return
}

func ReadVCardMessageNode(message *waProto.Message) (displayName string, contacts []string) {
	reg := regexp.MustCompile(`\d{7,15}`)

	if message.ContactMessage != nil {
		displayName = message.ContactMessage.GetDisplayName()

		results := reg.FindAllStringSubmatch(message.ContactMessage.GetVcard(), -1)
		if len(results) > 0 {
			contacts = append(contacts, results[0][1])
		}
	} else if message.ContactsArrayMessage != nil {
		displayName = message.ContactsArrayMessage.GetDisplayName()

		cards := message.ContactsArrayMessage.GetContacts()
		for _, card := range cards {
			results := reg.FindAllStringSubmatch(card.GetVcard(), -1)
			if len(results) > 0 {
				contacts = append(contacts, results[0][1])
			}
		}
	}

	return
}

func castPhoneNumber(number string) (international, telephone string, err error) {
	var imsi msisdn.IMSIParser
	imsi, err = msisdn.Parse(number, true)
	if err != nil {
		return
	}

	var phone *phonenumbers.PhoneNumber
	phone, err = phonenumbers.Parse(number, imsi.GetISO())
	if err != nil {
		return
	}

	international = phonenumbers.Format(phone, phonenumbers.INTERNATIONAL)
	telephone = fmt.Sprintf("%v%v", phone.GetCountryCode(), phone.GetNationalNumber())
	return
}
