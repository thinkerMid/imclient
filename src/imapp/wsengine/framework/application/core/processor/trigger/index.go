package trigger

import (
	"ws/framework/application/constant/binary"
	"ws/framework/application/container/abstract_interface"
)

type triggerState uint8

const (
	waitForTrigger triggerState = iota
	triggered
	disable
)

// Trigger 触发判断器
//
//	用于放行消息给处理器执行
type Trigger struct {
	tag       string
	status    triggerState
	automatic bool // false 则为依赖关键字触发，true 为自动触发
}

// Tag .
func (t *Trigger) Tag() string {
	return t.tag
}

// NoneTrigger .
func NoneTrigger() containerInterface.ITrigger {
	return &Trigger{status: disable}
}

// NewDefaultTrigger .
func NewDefaultTrigger(key string) containerInterface.ITrigger {
	return &Trigger{tag: key, status: waitForTrigger}
}

// NewAutomaticTrigger .
func NewAutomaticTrigger() containerInterface.ITrigger {
	return &Trigger{automatic: true, status: waitForTrigger}
}

// DefaultTrigger .
func (t *Trigger) DefaultTrigger() bool {
	return t.automatic
}

// WaitActive .
func (t *Trigger) WaitActive() bool {
	return t.status == waitForTrigger
}

// Disable .
func (t *Trigger) Disable() {
	t.status = disable
}

// Reset .
func (t *Trigger) Reset() {
	// 非自动才可以重置回开启状态
	if !t.automatic {
		t.status = waitForTrigger
	}
}

// ContentedCondition .
func (t *Trigger) ContentedCondition(node *waBinary.Node) bool {
	if t.status == triggered {
		return true
	}

	if t.tag != node.Tag {
		return false
	}

	t.status = triggered

	return true
}
