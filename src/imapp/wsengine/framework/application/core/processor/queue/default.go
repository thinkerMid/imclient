package actionQueue

import (
	"container/list"
	"ws/framework/application/container/abstract_interface"
	functionTools "ws/framework/utils/function_tools"
)

func convertLinkList(actionQueue []containerInterface.IAction) *list.List {
	linkList := list.New()

	var i uint8

	for _, action := range actionQueue {
		i++

		action.SetActionIndex(i)
		action.SetActionName(functionTools.ReflectValueTypeName(action))

		linkList.PushBack(action)
	}

	return linkList
}

// Default .
//
//	会一直弹出action给外部执行
//	如果遇到 error 会中断 Pop
type Default struct {
	actionQueue *list.List
	size        int
	whenError   bool
}

// NewDefault .
func NewDefault(actionQueue []containerInterface.IAction) containerInterface.IActionQueue {
	return &Default{
		actionQueue: convertLinkList(actionQueue),
		size:        len(actionQueue),
	}
}

// Name .
func (m *Default) Name() string {
	return "default"
}

// Size .
func (m *Default) Size() int {
	return m.size
}

// Empty .
func (m *Default) Empty() bool {
	return m.actionQueue.Len() == 0
}

// OnError .
func (m *Default) OnError(err error) {
	m.whenError = true
}

// Current .
func (m *Default) Current() containerInterface.IAction {
	element := m.actionQueue.Front()
	if element != nil {
		return element.Value.(containerInterface.IAction)
	}

	return nil
}

// Pop .
func (m *Default) Pop() containerInterface.IAction {
	m.actionQueue.Remove(m.actionQueue.Front())

	if m.whenError || m.actionQueue.Len() == 0 {
		// 清空
		m.actionQueue.Init()
		return nil
	}

	return m.Current()
}

// Dump .
func (m *Default) Dump() (actionNames string) {
	if m.actionQueue.Len() == 0 {
		return
	}

	nextElement := m.actionQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		actionNames += currentElement.Value.(containerInterface.IAction).ActionName()
		actionNames += ","
	}

	actionNames = actionNames[:len(actionNames)-1]

	return
}
