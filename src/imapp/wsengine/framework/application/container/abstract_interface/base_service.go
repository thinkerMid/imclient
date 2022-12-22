package containerInterface

import (
	"go.uber.org/zap"
	"ws/framework/application/constant/types"
)

// IService .
type IService interface {
	Init()
	SetJID(types.JID)
	SetAppIocContainer(IAppIocContainer)
	SetLogger(*zap.SugaredLogger)
	OnApplicationStart()
	OnApplicationResume()
	OnApplicationExit()
}

// IDataStorageService .
type IDataStorageService interface {
	IService
	CleanupAllData()
	OnJIDChangeWhenRegisterSuccess(newJID types.JID)
}

// BaseService IOC初始化的时候会调用函数进行赋值
type BaseService struct {
	AppIocContainer IAppIocContainer   // 容器
	JID             types.JID          // 上下文对象ID
	Logger          *zap.SugaredLogger // 日志打印
}

// Init 初始化
func (b *BaseService) Init() {}

// SetAppIocContainer .
func (b *BaseService) SetAppIocContainer(ioc IAppIocContainer) {
	b.AppIocContainer = ioc
}

// SetLogger .
func (b *BaseService) SetLogger(logger *zap.SugaredLogger) {
	b.Logger = logger
}

// SetJID .
func (b *BaseService) SetJID(jid types.JID) {
	b.JID = jid
}

// CleanupAllData 清空存储
func (b *BaseService) CleanupAllData() {}

// OnApplicationStart 实例退出
func (b *BaseService) OnApplicationStart() {}

// OnApplicationResume 实例退出
func (b *BaseService) OnApplicationResume() {}

// OnApplicationExit 实例退出
func (b *BaseService) OnApplicationExit() {}

// OnJIDChangeWhenRegisterSuccess JID不一致
func (b *BaseService) OnJIDChangeWhenRegisterSuccess(newJID types.JID) {}
