package businessService

import (
	"google.golang.org/protobuf/proto"
	"ws/framework/application/constant/binary/proto"
	containerInterface "ws/framework/application/container/abstract_interface"
	businessDB "ws/framework/application/data_storage/business/database"
	"ws/framework/application/libsignal/ecc"
	"ws/framework/plugin/database"
	databaseTools "ws/framework/plugin/database/database_tools"
	"ws/framework/utils"
	functionTools "ws/framework/utils/function_tools"
)

var _ containerInterface.IBusinessService = &Business{}

// Business .
type Business struct {
	containerInterface.BaseService

	context *businessDB.BusinessProfile
}

// Create .
func (d *Business) Create() (*businessDB.BusinessProfile, error) {
	context := businessDB.BusinessProfile{JID: d.JID.User}

	_, err := databaseTools.Create(database.MasterDB(), &context)
	if err != nil {
		d.Logger.Error(err)
		return nil, err
	}

	d.context = &context

	return &context, nil
}

// Context .
func (d *Business) Context() *businessDB.BusinessProfile {
	if d.context != nil {
		return d.context
	}

	context := businessDB.BusinessProfile{JID: d.JID.User}

	err := databaseTools.FindByPrimaryKey(database.MasterDB(), &context)
	if err != nil {
		d.Logger.Error(err)
		return nil
	}

	d.context = &context

	return &context
}

// GenerateBusinessVerifiedName .
func (d *Business) GenerateBusinessVerifiedName(appendPushName bool) (buff []byte) {
	device := d.AppIocContainer.ResolveDeviceService().Context()

	var verifiedName string
	if appendPushName && len(device.PushName) > 0 {
		verifiedName = device.PushName
	}

	// 随机17位整型作为序号
	serial := utils.RandInt64(10000000000000000, 99999999999999999)

	vnDetail := waProto.VerifiedNameDetails{
		Serial:       proto.Uint64(uint64(serial)),
		Issuer:       proto.String("smb:wa"),
		VerifiedName: proto.String(verifiedName),
	}

	buff, _ = proto.Marshal(&vnDetail)

	priKey := functionTools.SliceTo32SizeArray(device.IdentityKey)
	signature := ecc.CalculateSignature(ecc.NewDjbECPrivateKey(priKey), buff)

	var vnc = waProto.VerifiedNameCertificate{
		Details:   buff,
		Signature: signature[:],
	}

	buff, _ = proto.Marshal(&vnc)

	// update profile
	d.ContextExecute(func(context *businessDB.BusinessProfile) {
		context.UpdateVNameSerial(serial)
	})
	return
}

// ContextExecute .
func (d *Business) ContextExecute(f func(context *businessDB.BusinessProfile)) {
	context := d.Context()
	if context == nil {
		return
	}

	f(context)

	_, err := databaseTools.Save(database.MasterDB(), context)
	if err != nil {
		d.Logger.Error(err)
	}
}

// CleanupAllData .
func (d *Business) CleanupAllData() {
	context := businessDB.BusinessProfile{JID: d.JID.User}

	_, _ = databaseTools.DeleteByPrimaryKey(database.MasterDB(), &context)
}
