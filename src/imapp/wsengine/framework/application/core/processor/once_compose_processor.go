package processor

import (
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor/queue"
)

// NewOnceComposeProcessor .
func NewOnceComposeProcessor(actions [][]containerInterface.IAction, sets ...SetConfigFn) containerInterface.IProcessor {
	p := newLogicProcessor()

	for _, set := range sets {
		set(&p)
	}

	p.SetActionQueue(actionQueue.NewMultiple(actions))

	return &p
}
