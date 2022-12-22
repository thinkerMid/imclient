package eventService

import (
	"strconv"
	"time"
	"ws/framework/application/container/abstract_interface"
	eventDB "ws/framework/application/data_storage/event/database"
	"ws/framework/plugin/database"
	"ws/framework/plugin/database/database_tools"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

const (
	write int32 = iota
	standbyWrite

	// 日志缓存渠道id
	channel0 byte = 0
	channel1 byte = 1
	channel2 byte = 2

	cacheSize int = 2048
)

var _ containerInterface.IEventCache = &Channel0EventCache{}

// Channel0EventCache .
type Channel0EventCache struct {
	containerInterface.BaseService

	serialNumber  int64 // 日志批次序号
	autoIncrement int32 // 事件序号
	channelID     byte

	dirtyCode int32

	cacheEventBuffer eventSerialize.IEventBuffer
}

// Init .
func (e *Channel0EventCache) Init() {
	e.serialNumber = time.Now().UnixNano()
	e.channelID = channel0
	e.dirtyCode = standbyWrite
}

// ----------------------------------------------------------------------------

func (e *Channel0EventCache) mallocEventBuffer() eventSerialize.IEventBuffer {
	if e.cacheEventBuffer == nil {
		e.cacheEventBuffer = eventSerialize.AcquireEventBuffer()
	}

	data := e.cacheEventBuffer.Byte()
	if len(data) > cacheSize {
		e.store(data)

		eventSerialize.ReleaseEventBuffer(e.cacheEventBuffer)
		e.cacheEventBuffer = eventSerialize.AcquireEventBuffer()
	}

	return e.cacheEventBuffer
}

// AddEvent .
func (e *Channel0EventCache) AddEvent(event containerInterface.WaEvent) {
	e.dirtyCode = write

	event.Serialize(e.mallocEventBuffer())
}

// ----------------------------------------------------------------------------

// ClearLog 清除所有事件
func (e *Channel0EventCache) ClearLog() {
	err := eventDB.Delete(database.MasterDB(), e.JID.User, e.serialNumber, e.channelID)
	if err != nil {
		e.Logger.Error(err)
	}
}

// ClearNotSentYetLog 清除上次没发的事件
func (e *Channel0EventCache) ClearNotSentYetLog() {
	err := eventDB.DeleteLast(database.MasterDB(), e.JID.User, e.serialNumber, e.channelID)
	if err != nil {
		e.Logger.Error(err)
	}
}

// ----------------------------------------------------------------------------

// PackNotSentYetLog 打包上次没发的事件
func (e *Channel0EventCache) PackNotSentYetLog(sendEventCount int32, buffer eventSerialize.IEventBuffer) {
	e.generateUploadHeader(sendEventCount, buffer)

	e.fillLogRecord(false, buffer)
}

// PackBuffer 打包所有事件
func (e *Channel0EventCache) PackBuffer(sendEventCount int32, buffer eventSerialize.IEventBuffer) {
	e.generateUploadHeader(sendEventCount, buffer)

	e.fillLogRecord(true, buffer)

	e.serialNumber = time.Now().UnixNano()
}

// CacheBufferItem 已缓存的事件总数
func (e *Channel0EventCache) CacheBufferItem() int64 {
	if e.cacheEventBuffer != nil {
		if len(e.cacheEventBuffer.Byte()) > 0 {
			return 1
		}
	}

	count, err := eventDB.Count(database.MasterDB(), e.JID.User, e.serialNumber, e.channelID)
	if err != nil {
		e.Logger.Error(err)
	}

	return count
}

// ----------------------------------------------------------------------------

// ResetAddLogState 重置日志写状态
func (e *Channel0EventCache) ResetAddLogState() {
	if e.dirtyCode == write {
		e.dirtyCode = standbyWrite

		e.mallocEventBuffer().Common().SerializeNumber(8, float64(utils.GetCurTime()))
	}
}

// ----------------------------------------------------------------------------

// 从数据库取出日志内容
//
//	takeAll 为false则用于取serialNumber之前的日志
//	takeAll 为true则用于取serialNumber之前的和当前日志
func (e *Channel0EventCache) fillLogRecord(takeAll bool, buffer eventSerialize.IEventBuffer) {
	var allCacheBuffer []eventDB.EventBuffer
	var err error

	if takeAll {
		allCacheBuffer, err = eventDB.FindAll(database.MasterDB(), e.JID.User, e.serialNumber, e.channelID)
	} else {
		allCacheBuffer, err = eventDB.FindLast(database.MasterDB(), e.JID.User, e.serialNumber, e.channelID)
	}

	if err != nil {
		e.Logger.Error(err)
	}

	// 读出数据库的
	for _, b := range allCacheBuffer {
		buffer.Write(b.EventLog)
	}

	// 取所有的时候 看缓冲内有没有
	if takeAll && e.cacheEventBuffer != nil {
		data := e.cacheEventBuffer.Byte()
		if len(data) > 0 {
			buffer.Write(data)

			// 把缓冲存到数据库
			e.store(data)
		}

		// 释放
		eventSerialize.ReleaseEventBuffer(e.cacheEventBuffer)
		e.cacheEventBuffer = nil
	}
}

// 事件储存
func (e *Channel0EventCache) store(buff []byte) {
	e.autoIncrement++

	data := eventDB.EventInfo{}
	data.JID = e.JID.User
	data.SerialNumber = e.serialNumber
	data.AutoIncrement = e.autoIncrement
	data.ChannelID = e.channelID
	data.EventLog = buff

	_, err := databaseTools.Create(database.MasterDB(), &data)
	if err != nil {
		e.Logger.Error(err)
	}
}

// FlushEventCache .
func (e *Channel0EventCache) FlushEventCache() {
	// 缓冲内有没有
	if e.cacheEventBuffer != nil {
		data := e.cacheEventBuffer.Byte()
		if len(data) > 0 {
			// 把缓冲存到数据库
			e.store(data)
		}

		// 释放
		eventSerialize.ReleaseEventBuffer(e.cacheEventBuffer)
		e.cacheEventBuffer = nil
	}
}

// CleanupAllData .
func (e *Channel0EventCache) CleanupAllData() {
	if e.cacheEventBuffer != nil {
		eventSerialize.ReleaseEventBuffer(e.cacheEventBuffer)
		e.cacheEventBuffer = nil
	}

	e.ClearLog()
}

// ----------------------------------------------------------------------------

const (
	wam   int32 = 0x54D4157 // 固定值
	state byte  = 1
)

func (e *Channel0EventCache) generateUploadHeader(sendEventCount int32, buffer eventSerialize.IEventBuffer) {
	// 头部
	{
		// WAM + 0x05 + state + times + channelID
		// 57414d 05 01 3601 00
		buffer.WriteLittleEndianInt32(wam, false)
		buffer.WriteByte(state)
		buffer.WriteLittleEndianInt16(int16(sendEventCount), false)
		buffer.WriteByte(e.channelID)
	}

	// 渠道2没有基础包
	if e.channelID == channel2 {
		return
	}

	device := e.AppIocContainer.ResolveDeviceService().Context()
	configuration := e.AppIocContainer.ResolveWhatsappConfiguration()
	login := e.AppIocContainer.ResolveMemoryCache().AccountLoginData()

	// 基础包
	{
		// 设置基础包的上下文
		buffer.Common()

		if len(device.Mnc) != 0 {
			mnc, err := strconv.Atoi(device.Mnc)
			if err != nil {
				buffer.SerializeNumber(configuration.CommonSerializeCode[0], float64(mnc))
			}
		}

		if len(device.Mcc) != 0 {
			mcc, err := strconv.Atoi(device.Mcc)
			if err != nil {
				buffer.SerializeNumber(configuration.CommonSerializeCode[0], float64(mcc))
			}
		}

		buffer.SerializeNumber(configuration.CommonSerializeCode[2], 1.000000)
		buffer.SerializeString(configuration.CommonSerializeCode[3], device.Device)               //设备机型  iPhone 6s
		buffer.SerializeString(configuration.CommonSerializeCode[4], device.OsVersion)            //设备版本  13.6
		buffer.SerializeString(configuration.CommonSerializeCode[5], configuration.VersionString) //客户端版本号  已确认
		buffer.SerializeNumber(configuration.CommonSerializeCode[6], 0.000000)                    //固定的
		buffer.SerializeNumber(configuration.CommonSerializeCode[7], 0.000000)                    //是否是wifi 0不是 1是
		buffer.SerializeNumber(configuration.CommonSerializeCode[8], float64(utils.GetCurTime())) //时间戳
		buffer.SerializeNumber(configuration.CommonSerializeCode[9], 111.000000)                  //网络状态 111:lte
		buffer.SerializeNumber(configuration.CommonSerializeCode[10], 1.000000)                   // iphone process 固定1
		buffer.SerializeNumber(configuration.CommonSerializeCode[13], 4.000000)                   // 固定4
		buffer.SerializeNumber(configuration.CommonSerializeCode[14], 2.000000)                   // distributionChannel 固定2

		//if len(login.ABKey) != 0 {
		//	serializeCommonField(cache, 16, string(login.ABKey)) 	// 真机新版本不赋值
		//}

		if len(login.LastKnowDataCenter) != 0 {
			buffer.SerializeString(configuration.CommonSerializeCode[18], login.LastKnowDataCenter)
		}

		if len(login.ABKey2) != 0 {
			buffer.SerializeString(configuration.CommonSerializeCode[20], login.ABKey2)
		}

		//22 ABExposureKey
		if len(login.ABExposureKey) != 0 {
			buffer.SerializeString(configuration.CommonSerializeCode[22], string(login.ABExposureKey))
		}

		if sendEventCount == 0 {
			if len(device.PrivateStatsId) != 0 {
				buffer.SerializeString(configuration.CommonSerializeCode[23], device.PrivateStatsId) //注册时发，正常登录不发
			}
		}

		buffer.SerializeNumber(configuration.CommonSerializeCode[24], 1.000000)       //设备环境+越狱检测 (这里的值=NOT(flag) AND 1)
		buffer.SerializeString(configuration.CommonSerializeCode[25], device.Country) //国家
		//buffer.SerializeNumber(26, 0.000000)               // mdCompanionRegOptIn 固定0 ??? 真机不发
		buffer.SerializeNumber(configuration.CommonSerializeCode[27], 1.000000)                      // WAPreferences 的 mdCompanionRegOptIn
		buffer.SerializeString(configuration.CommonSerializeCode[28], device.BuildNumber)            //iOSBuildNumber
		buffer.SerializeString(configuration.CommonSerializeCode[29], configuration.BuildSDKVersion) //新增 sdk版本
	}
}
