package application

import (
	"container/list"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor/watch_dog"
	"ws/framework/plugin/lightning"
)

// IApplication .
type IApplication interface {
	Container() containerInterface.IAppIocContainer
	Start() error
	Exit()
}

// App .
type App struct {
	container containerInterface.IAppIocContainer // 容器
	logger    *zap.SugaredLogger                  // 日志

	status uint32 // 状态: 默认standby

	idCounter  int
	sidCounter int

	messageProcessorId    uint32     // 消息处理器的自增PID
	messageProcessorQueue *list.List // 消息处理器链表

	globalResultProcessorQueue []containerInterface.GlobalResultProcessor // 全局消息结果监听队列

	workerPool          *lightning.WorkPool // 消息处理循环
	timeSchedule        *time.Timer         // 取消定时器的句柄
	shutdownErrorRemark shutdownErrorRemark // 流程终止之前的异常

	watchDog *watchDog.WatchDog
}

func (c *App) init() {
	// 事件循环池
	c.workerPool = lightning.NewAsyncPool(2<<7, c.eventLoop)
	// 消息处理器队列
	c.messageProcessorQueue = list.New()
	// 执行流程监控
	c.watchDog = watchDog.New(c.logger)
	// 定时器
	c.timeSchedule = time.AfterFunc(time.Second, c.timeTick)
	c.timeSchedule.Stop()
}

// Container .
func (c *App) Container() containerInterface.IAppIocContainer {
	return c.container
}

// Start 发起连接->握手->启动成功
func (c *App) Start() (err error) {
	c.init()

	err = c.container.ResolveConnection().Connect()
	if err != nil {
		return
	}

	c.logger.Info("start")
	c.container.OnStart()

	c.changeStatus(activate)

	// 开始监听连接回包
	c.container.ResolveConnection().SetEventListener(c)
	// 定时器
	c.timeSchedule.Reset(time.Second)
	return
}

// 重连成功
func (c *App) resume() {
	c.logger.Info("resume")
	c.container.OnResume()

	// 开始监听连接回包
	c.container.ResolveConnection().SetEventListener(c)
	// 定时器
	c.timeSchedule.Reset(time.Second)
}

// Exit .
// 外部场景调用，进来的是非异步事件的协程
func (c *App) Exit() {
	c.changeStatus(exit)

	c.workerPool.Invoke(newSignalEvent(destroy), false)
	c.workerPool.Await()

	c.logger.Info("exit")
	c.container.OnExit()
}

func (c *App) changeStatus(v uint32) {
	status := atomic.LoadUint32(&c.status)

	if status == v {
		return
	}

	atomic.CompareAndSwapUint32(&c.status, status, v)
}
