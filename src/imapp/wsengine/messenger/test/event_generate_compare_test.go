package test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	waProto "ws/framework/application/constant/binary/proto"
	appContainer "ws/framework/application/container"
	"ws/framework/application/container/abstract_interface"
	wam2 "ws/framework/application/core/wam"
	"ws/framework/application/core/wam/events"
	"ws/framework/application/data_storage/cache"
	"ws/framework/application/data_storage/device/database"
	"ws/framework/engine/consts"
	"ws/framework/env"
	eventSerialize "ws/framework/plugin/event_serialize"
	"ws/framework/utils"
)

type Cache struct {
	Buffer eventSerialize.IEventBuffer

	A []byte
	B []byte
}

func (c *Cache) AddEvent(e containerInterface.WaEvent) {
	c.Buffer.Reset()
	e.Serialize(c.Buffer)

	c.A = append(c.A, c.Buffer.Byte()...)

	//c.Buffer.Println("wa2")
}

func (c *Cache) CacheBufferItem() int64 {
	return 0
}

func (c *Cache) PackBuffer(sendEventCount int32, buffer eventSerialize.IEventBuffer) {
}

func (c *Cache) ClearLog() {}

func (c *Cache) ResetAddLogState() {}

func (c *Cache) ClearNotSentYetLog() {}

func (c *Cache) PackNotSentYetLog(sendEventCount int32, buffer eventSerialize.IEventBuffer) {
}

func (c *Cache) FlushEventCache() {}

func (c *Cache) Cleanup() {}

// ----------------------------------------------------------------------------

type Device struct{}

func (d Device) Context() *deviceDB.Device {
	return &deviceDB.Device{}
}

func (d Device) Create(virtualDevice containerInterface.IVirtualDevice) (*deviceDB.Device, error) {
	//TODO implement me
	panic("implement me")
}

func (d Device) Import(device *deviceDB.Device) error {
	//TODO implement me
	panic("implement me")
}

func (d Device) GetClientPayload() *waProto.ClientPayload {
	//TODO implement me
	panic("implement me")
}

func (d Device) DeviceAgent() string {
	//TODO implement me
	panic("implement me")
}

func (d Device) PrivateStatsAgent() string {
	//TODO implement me
	panic("implement me")
}

func (d Device) ContextExecute(f func(*deviceDB.Device)) {
	//TODO implement me
	panic("implement me")
}

// ----------------------------------------------------------------------------

func TestGenerateCompare(t *testing.T) {
	c1 := &Cache{Buffer: eventSerialize.AcquireEventBuffer()}
	c2 := &Cache{Buffer: eventSerialize.AcquireEventBuffer()}
	d := &Device{}

	ioc := appContainer.NewAppIocContainer()
	ioc.Inject(appContainer.Channel0EventCache, c1)
	ioc.Inject(appContainer.Channel2EventCache, c2)
	ioc.Inject(appContainer.Device, d)

	type testCase struct {
		name      string
		executeFn func()
	}

	caseFn := []testCase{
		{
			"LogLogin",
			func() {
				wam2.LogManager().LogLogin(ioc)
			},
		},
		{
			"LogRegisterLaunch",
			func() {
				wam2.LogManager().LogRegisterLaunch(ioc)
			},
		},
		{
			"LogContactAdd:haveAvatar",
			func() {
				wam2.LogManager().LogContactAdd(ioc, true, false)
			},
		},
		{
			"LogContactAdd:noHaveAvatar",
			func() {
				wam2.LogManager().LogContactAdd(ioc, true, true)
			},
		},
		{
			"LogDeleteContact:haveAvatar",
			func() {
				wam2.LogManager().LogDeleteContact(ioc, true, false)
			},
		},
		{
			"LogDeleteContact:noHaveAvatar",
			func() {
				wam2.LogManager().LogDeleteContact(ioc, true, false)
			},
		},
		{
			"LogNotifyContactAvatar:1024size",
			func() {
				wam2.LogManager().LogNotifyContactAvatar(ioc, 1024)
			},
		},
		{
			"LogNotifyContactAvatar:0size",
			func() {
				wam2.LogManager().LogNotifyContactAvatar(ioc, 0)
			},
		},
		{
			"LogNotifyContactAddedOrDeleted",
			func() {
				wam2.LogManager().LogNotifyContactAddedOrDeleted(ioc)
			},
		},
		{
			"LogSendText:first",
			func() {
				wam2.LogManager().LogSendText(ioc, 0, true)
			},
		},
		{
			"LogSendText:second",
			func() {
				wam2.LogManager().LogSendText(ioc, 2, false)
			},
		},
		{
			"LogSendText:MediaImage",
			func() {
				wam2.LogManager().LogSendMedia(ioc, events.MediaImage, wam2.Media{Image: &wam2.Image{
					Width:     1024,
					Height:    1024,
					Size:      1024,
					FirstScan: 1024,
					LowScan:   1024,
					MidScan:   1024,
				}})
			},
		},
		{
			"LogSendText:MediaVideo",
			func() {
				wam2.LogManager().LogSendMedia(ioc, events.MediaVideo, wam2.Media{Video: &wam2.Video{
					Width:  1024,
					Height: 1024,
					Size:   1024,
				}})
			},
		},
		{
			"LogSendText:MediaVoice",
			func() {
				wam2.LogManager().LogSendMedia(ioc, events.MediaVoice, wam2.Media{Voice: &wam2.Voice{Size: 1024}})
			},
		},
		{
			"LogSendText:MediaVCard",
			func() {
				wam2.LogManager().LogSendMedia(ioc, events.MediaVCard, wam2.Media{})
			},
		},
		{
			"LogSessionNew",
			func() {
				wam2.LogManager().LogSessionNew(ioc)
			},
		},
	}

	for _, c := range caseFn {
		c1.A = make([]byte, 0)
		c1.B = make([]byte, 0)
		c1.Buffer.Reset()

		t.Run(c.name, func(t *testing.T) {
			c.executeFn()
		})
	}
}

func TestCommonSerialize(t *testing.T) {
	device := deviceDB.Device{
		Device:      "iPhone 6s",
		OsVersion:   "13.6",
		Country:     "MY",
		BuildNumber: "17G68",
	}

	buffer := eventSerialize.AcquireEventBuffer()
	c := memoryCacheService.MemoryCache{}
	c.AccountLoginData().LastKnowDataCenter = "frc"
	c.AccountLoginData().ABKey2 = "2ST,HY,hn,7j,6L,3F,3H,1J,n,R,I,M,3,7,M,6K,1u,V"

	login := c.AccountLoginData()

	t.Run("TestCommonSerialize", func(t *testing.T) {
		// 头部
		//{
		//	const (
		//		wam   int32 = 0x54D4157 // 固定值
		//		state byte  = 1
		//	)
		//	// WAM + 0x05 + state + times + channelID
		//	// 57414d 05 01 3601 00
		//	buffer.WriteLittleEndianInt32(wam, false)
		//	buffer.WriteByte(state)
		//	buffer.WriteLittleEndianInt16(int16(289), false) // TODO 为什么要转int16?
		//	buffer.WriteByte(0)
		//}

		// 设置基础包的上下文
		buffer.Common()

		if len(device.Mnc) != 0 {
			mnc, err := strconv.Atoi(device.Mnc)
			if err != nil {
				buffer.SerializeCommonNumber(0, float64(mnc))
			}
		}

		if len(device.Mcc) != 0 {
			mcc, err := strconv.Atoi(device.Mcc)
			if err != nil {
				buffer.SerializeCommonNumber(0, float64(mcc))
			}
		}

		buffer.SerializeCommonNumber(2, 1.000000)
		buffer.SerializeCommonString(3, device.Device)                                  //设备机型  iPhone 6s
		buffer.SerializeCommonString(4, device.OsVersion)                               //设备版本  13.6
		buffer.SerializeCommonString(5, env.NacosConfig.WsaConfig.GetWAVersionString()) //客户端版本号  已确认
		buffer.SerializeCommonNumber(6, 0.000000)                                       //固定的
		buffer.SerializeCommonNumber(7, 1.000000)                                       //是否是wifi 0不是 1是
		buffer.SerializeCommonNumber(8, float64(utils.GetCurTime()))                    //时间戳
		buffer.SerializeCommonNumber(9, 1.000000)                                       //网络状态 111:lte
		buffer.SerializeCommonNumber(10, 1.000000)                                      // iphone process 固定1
		buffer.SerializeCommonNumber(13, 4.000000)                                      // 固定4
		buffer.SerializeCommonNumber(14, 2.000000)                                      // distributionChannel 固定2

		//if len(login.ABKey) != 0 {
		//	serializeCommonField(cache, 16, string(login.ABKey)) 	// 真机新版本不赋值
		//}

		if len(login.LastKnowDataCenter) != 0 {
			buffer.SerializeCommonString(18, login.LastKnowDataCenter)
		}

		if len(login.ABKey2) != 0 {
			buffer.SerializeCommonString(20, login.ABKey2)
		}

		//22 ABExposureKey
		if len(login.ABExposureKey) != 0 {
			buffer.SerializeCommonString(22, "4876,5254,5818,3803,4675,5079,1088,5236,5160")
		}

		buffer.SerializeCommonNumber(24, 1.000000)               //设备环境+越狱检测 这里填1是未检测出异常  ???测试发现0
		buffer.SerializeCommonString(25, device.Country)         //国家
		buffer.SerializeCommonNumber(26, 0.000000)               // mdCompanionRegOptIn 固定0 ??? 真机不发
		buffer.SerializeCommonNumber(27, 1.000000)               // WAPreferences 的 mdCompanionRegOptIn
		buffer.SerializeCommonString(28, device.BuildNumber)     //iOSBuildNumber
		buffer.SerializeCommonString(29, consts.BuildSDKVersion) //新增 sdk版本

		buffer.Println("common")

		dst, _ := hex.DecodeString("200b800d096950686f6e65203673800f0431332e3680110a322e32322e32302e373510152017502fb42f4963206928830138790604387b060288eb0a036672638879112e3253542c48592c686e2c376a2c364c2c33462c33482c314a2c6e2c522c492c4d2c332c372c4d2c364b2c31752c5688a5132c353831382c343837362c333830332c353037392c353233362c353235342c313038382c353136302c34363735186b1888b11a024d5928a71c88911e0531374736388879240431352e3529e8047601250681f577c7394129d006720261796b9925ef063f7601008050414fcb394129ee0122058202134d61696e546872656164426c6f636b656436303206072604")

		fmt.Println("[common] buffer hex:", "200b800d096950686f6e65203673800f0431332e3680110a322e32322e32302e373510152017502fb42f4963206928830138790604387b060288eb0a036672638879112e3253542c48592c686e2c376a2c364c2c33462c33482c314a2c6e2c522c492c4d2c332c372c4d2c364b2c31752c5688a5132c353831382c343837362c333830332c353037392c353233362c353235342c313038382c353136302c34363735186b1888b11a024d5928a71c88911e0531374736388879240431352e3529e8047601250681f577c7394129d006720261796b9925ef063f7601008050414fcb394129ee0122058202134d61696e546872656164426c6f636b656436303206072604")

		t.Log(buffer.Byte())
		t.Log(dst)
		if !bytes.Equal(buffer.Byte(), dst) {
			t.Error("not compare")
		}
	})
}
