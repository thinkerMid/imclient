package privateChatCommon

import (
	"github.com/nyaruka/phonenumbers"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
	messageResultType "ws/framework/application/core/result/constant"
	"ws/framework/lib/msisdn"
)

// VCardPhoneNumberCheck .
type VCardPhoneNumberCheck struct {
	processor.BaseAction
	Contacts []string
}

// Start .
func (m *VCardPhoneNumberCheck) Start(_ containerInterface.IMessageContext, next containerInterface.NextActionFn) (err error) {
	var imsi msisdn.IMSIParser

	for i := range m.Contacts {
		imsi, err = msisdn.Parse(m.Contacts[i], true)
		if err != nil {
			return
		}

		_, err = phonenumbers.Parse(m.Contacts[i], imsi.GetISO())
		if err != nil {
			return
		}
	}

	next()
	return
}

// Receive .
func (m *VCardPhoneNumberCheck) Receive(_ containerInterface.IMessageContext, _ containerInterface.NextActionFn) error {
	return nil
}

// Error .
func (m *VCardPhoneNumberCheck) Error(context containerInterface.IMessageContext, err error) {
	context.AppendResult(containerInterface.MessageResult{
		ResultType: messageResultType.VCardCheck,
		Error:      err,
	})
}
