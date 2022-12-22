package application

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"ws/framework/application/constant"
	"ws/framework/application/constant/binary"
	"ws/framework/application/constant/message"
	"ws/framework/plugin/logger"
	"ws/framework/utils"
)

// shutdownErrorRemark 流程终止之前的异常
type shutdownErrorRemark struct {
	accountStateErr error // 账号异常 通常为最直接体现账号状态401，403
	streamErr       error // 通信异常 预防未能知道账号异常，导致的无异常现象，捕获较为低级的异常作为补偿
}

func (s *shutdownErrorRemark) accountError(err error) {
	s.accountStateErr = err
}

func (s *shutdownErrorRemark) streamError(err error) {
	s.streamErr = err
}

func (s *shutdownErrorRemark) hasError() bool {
	return s.accountStateErr != nil || s.streamErr != nil
}

func (s *shutdownErrorRemark) thrown() error {
	if s.accountStateErr != nil {
		return s.accountStateErr
	}

	return s.streamErr
}

// GenerateRequestID .
func (c *App) GenerateRequestID() string {
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	c.idCounter++
	requestId := curTime + "-" + strconv.Itoa(c.idCounter)
	return requestId
}

// GenerateSID .
func (c *App) GenerateSID() string {
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	randNum := strconv.FormatInt(utils.RandInt64(1, 0xffffffff), 10)
	c.sidCounter++
	sid := curTime + ("-" + randNum + "-" + strconv.Itoa(c.sidCounter))
	return sid
}

// SendIQ .
func (c *App) SendIQ(query message.InfoQuery) (string, error) {
	if len(query.ID) == 0 {
		return "", fmt.Errorf("InfoQuery `ID` must not be empty")
	}

	attrs := waBinary.Attrs{
		"id":    query.ID,
		"xmlns": query.Namespace,
		"type":  query.Type,
	}

	if !query.To.IsEmpty() {
		attrs["to"] = query.To
	}

	if !query.Target.IsEmpty() {
		attrs["target"] = query.Target
	}

	node := waBinary.Node{
		Tag:     "iq",
		Attrs:   attrs,
		Content: query.Content,
	}

	return query.ID, c.sendPayload(node)
}

// SendNode .
func (c *App) SendNode(node waBinary.Node) (string, error) {
	return node.ID(), c.sendPayload(node)
}

func (c *App) sendPayload(node waBinary.Node) error {
	if c.shutdownErrorRemark.hasError() {
		return c.shutdownErrorRemark.thrown()
	}

	if logger.EnabledDebug() {
		c.logger.Debug("[C] ", node.XMLString())
	}

	return c.container.ResolveConnection().Write(node)
}

// processMessage .
func (c *App) processMessage(node *waBinary.Node) {
	hasErr := c.decodeMessageError(node)

	if hasErr {
		c.logger.Info("[S] ", node.XMLString())
	} else if logger.EnabledDebug() {
		c.logger.Debug("[S] ", node.XMLString())
	}

	c.invokeMessage(node, false)
}

func (c *App) decodeMessageError(node *waBinary.Node) bool {
	// 捕获这些tag 作为shutdown的异常
	switch node.Tag {
	// 登录失败
	case message.Failure:
		c.shutdownErrorRemark.accountError(constant.NewIgnoreError(node.AttrGetter().String("reason")))

		// 直接关闭连接
		c.container.ResolveConnection().Close()
	// 通信异常
	case message.StreamError:
		code := node.AttrGetter().String("code")

		if len(code) > 0 {
			c.shutdownErrorRemark.streamError(constant.NewIgnoreError(code))
		} else {
			// 是否属于抢占登录触发的异常
			conflict, ok := node.GetOptionalChildByTag("conflict")
			if ok {
				if conflict.AttrGetter().String("type") == "replaced" {
					// 改变状态 避免双方无限重连
					c.changeStatus(exit)
				}
			}

			c.shutdownErrorRemark.streamError(errors.New(node.XMLString()))
		}
	// 通信关闭
	case message.StreamEnd:
		// 如果没收到过异常才设置 因为这个streamend是最后来的协议包 如果前面有异常可以捕获的就不覆盖了
		if c.shutdownErrorRemark.streamErr == nil {
			c.shutdownErrorRemark.streamError(constant.NewIgnoreError(message.StreamEnd))
		}

		// 直接关闭连接
		c.container.ResolveConnection().Close()
	default:
		return false
	}

	return true
}
