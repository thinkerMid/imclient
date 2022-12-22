package actionQueue

import (
	"container/list"
	"ws/framework/application/container/abstract_interface"
)

// Multiple .
type Multiple struct {
	actionQueueGroup []*list.List
	whenError        bool
	size             int
}

// NewMultiple .
func NewMultiple(actionQueueGroup [][]containerInterface.IAction) containerInterface.IActionQueue {
	linkLists := make([]*list.List, len(actionQueueGroup))

	for i := range actionQueueGroup {
		linkLists[i] = convertLinkList(actionQueueGroup[i])
	}

	return &Multiple{
		actionQueueGroup: linkLists,
		size:             linkLists[0].Len(),
	}
}

// Name .
func (m *Multiple) Name() string {
	return "multiple"
}

// Size .
func (m *Multiple) Size() int {
	return m.size
}

// Empty .
func (m *Multiple) Empty() bool {
	return len(m.actionQueueGroup) == 0
}

// OnError .
func (m *Multiple) OnError(err error) {
	m.whenError = true
}

// Current .
func (m *Multiple) Current() containerInterface.IAction {
	if len(m.actionQueueGroup) == 0 {
		return nil
	}

	element := m.actionQueueGroup[0].Front()
	if element != nil {
		return element.Value.(containerInterface.IAction)
	}

	return nil
}

// Pop .
func (m *Multiple) Pop() containerInterface.IAction {
	if len(m.actionQueueGroup) == 0 {
		return nil
	}

	nowLinkList := m.actionQueueGroup[0]
	nowLinkList.Remove(nowLinkList.Front())

	// 还可以继续执行 && 没有出现异常
	if nowLinkList.Len() > 0 && !m.whenError {
		return m.Current()
	}

	// 重置
	m.whenError = false
	nowLinkList.Init()

	if len(m.actionQueueGroup) > 1 {
		m.actionQueueGroup = m.actionQueueGroup[1:]
		m.size = m.actionQueueGroup[0].Len()

		return m.Current()
	}

	m.actionQueueGroup = nil

	return nil
}

// Dump .
func (m *Multiple) Dump() (actionNames string) {
	if len(m.actionQueueGroup) == 0 {
		return
	}

	nextElement := m.actionQueueGroup[0].Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		actionNames += currentElement.Value.(containerInterface.IAction).ActionName()
		actionNames += ","
	}

	actionNames = actionNames[:len(actionNames)-1]

	return
}
