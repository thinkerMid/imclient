package actionQueue

import "ws/framework/application/container/abstract_interface"

// IgnoreError .
//
//	继承了 Default 的特性
//	不同点在于出现了异常不会有任何影响，会一直 Pop
type IgnoreError struct {
	Default
}

// NewIgnoreError .
func NewIgnoreError(actionQueue []containerInterface.IAction) containerInterface.IActionQueue {
	q := IgnoreError{}
	q.actionQueue = convertLinkList(actionQueue)
	q.size = len(actionQueue)

	return &q
}

// Name .
func (m *IgnoreError) Name() string {
	return "ignoreError"
}

// OnError .
func (m *IgnoreError) OnError(err error) {}
