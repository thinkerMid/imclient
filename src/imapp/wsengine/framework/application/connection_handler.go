package application

import (
	//"ws/client/whatsapp/client/im/core"
	"ws/framework/application/constant/binary"
)

// reconnect 重连
//
//	这里作为Client生命周期的最后一道
//	 连接关闭会尝试重连，如果连接失败，得到连接异常的结果，如果成功，则重新初始化
//	 只要重连失败了，都会将shutdownError的异常等级作为最终异常抛到action中
func (c *App) reconnect() error {
	c.cleanupBackgroundProcessor()

	err := c.container.ResolveConnection().Connect()
	if err != nil {
		c.logger.Warnf("reconnect failure. err: %s", err)
		// 如果没有异常 才把这个异常赋值上去 免得连接失败把上次较为准确的异常给覆盖了
		if c.shutdownErrorRemark.streamErr == nil {
			c.shutdownErrorRemark.streamError(err)
		}
		return err
	}

	// 连上了代表可以知道账号异常 可以把通信异常清空
	c.shutdownErrorRemark.streamErr = nil

	return nil
}

// OnConnectionClose .
// 内部场景调用，进来的是网络库的协程
func (c *App) OnConnectionClose() {
	c.timeSchedule.Stop()
	c.workerPool.Invoke(newSignalEvent(disconnect), false)
}

// OnResponse .
// 内部场景调用，进来的是网络库的协程
func (c *App) OnResponse(node *waBinary.Node) {
	c.workerPool.Invoke(newMessageEvent(node), false)
}
