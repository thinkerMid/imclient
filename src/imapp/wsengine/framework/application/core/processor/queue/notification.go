package actionQueue

import (
	"container/list"
	"ws/framework/application/container/abstract_interface"
	functionTools "ws/framework/utils/function_tools"
)

// NotificationQueue .
type NotificationQueue struct {
	actionQueue *list.List

	size    int
	aborted bool

	makeActionQueueFn containerInterface.MakeNotificationFn
}

// NewNotification .
func NewNotification(makeActionQueueFn containerInterface.MakeNotificationFn) containerInterface.INotificationQueue {
	actionQueue := makeActionQueueFn()

	linkList := list.New()
	for _, action := range actionQueue {
		linkList.PushBack(action)
	}

	return &NotificationQueue{
		actionQueue:       linkList,
		makeActionQueueFn: makeActionQueueFn,
		size:              linkList.Len(),
	}
}

// Reload .
func (m *NotificationQueue) Reload() {
	m.actionQueue.Init()

	actionQueue := m.makeActionQueueFn()
	for _, action := range actionQueue {
		m.actionQueue.PushBack(action)
	}

	m.size = m.actionQueue.Len()
}

// Current .
func (m *NotificationQueue) Current() containerInterface.INotification {
	element := m.actionQueue.Front()
	if element != nil {
		return element.Value.(containerInterface.INotification)
	}

	return nil
}

// Pop .
func (m *NotificationQueue) Pop() containerInterface.INotification {
	if m.actionQueue.Len() > 0 {
		element := m.actionQueue.Front()
		m.actionQueue.Remove(element)

		return element.Value.(containerInterface.INotification)
	}

	return nil
}

// Dump .
func (m *NotificationQueue) Dump() (actionNames string) {
	if m.actionQueue.Len() == 0 {
		return
	}

	nextElement := m.actionQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		actionNames += functionTools.ReflectValueTypeName(currentElement.Value)
		actionNames += ","
	}

	actionNames = actionNames[:len(actionNames)-1]

	return
}
