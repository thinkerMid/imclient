package processor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor/queue"
)

// NewOnceIgnoreErrorProcessor .
func NewOnceIgnoreErrorProcessor(actions []containerInterface.IAction, sets ...SetConfigFn) containerInterface.IProcessor {
	p := newLogicProcessor()

	for _, set := range sets {
		set(&p)
	}

	p.SetActionQueue(actionQueue.NewIgnoreError(actions))

	return &p
}
