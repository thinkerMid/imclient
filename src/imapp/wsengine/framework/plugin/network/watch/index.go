package networkWatch

import (
	"container/list"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"ws/framework/plugin/lightning"
)

type event struct {
	code    uint8
	content interface{}
}

type target struct {
	id           uint32
	c            net.Conn
	nowBlockTime uint8
	maxBlockTime uint8
}

// ConnectionWatch .
type ConnectionWatch struct {
	scheduler *time.Timer
	worker    *lightning.WorkPool

	incrementID uint32

	watchQueue *list.List
}

var watchInit sync.Once
var watcher *ConnectionWatch

// Instance .
func Instance() *ConnectionWatch {
	watchInit.Do(func() {
		watcher = &ConnectionWatch{
			watchQueue: list.New(),
		}

		watcher.scheduler = time.AfterFunc(time.Second, watcher.wakeUp)
		watcher.worker = lightning.New(1, 2<<9, watcher.loop)

		watcher.scheduler.Stop()
	})

	return watcher
}

func (w *ConnectionWatch) loop(i interface{}) {
	e := i.(event)
	switch e.code {
	case 1:
		w.queueAdd(e.content.(*target))
	case 2:
		w.queueRemove(e.content.(*target))
	case 3:
		w.check()
	}
}

func (w *ConnectionWatch) wakeUp() {
	w.worker.Invoke(event{code: 3}, false)
}

func (w *ConnectionWatch) check() {
	nextElement := w.watchQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		t := currentElement.Value.(*target)

		t.nowBlockTime += 1

		if t.nowBlockTime < t.maxBlockTime {
			continue
		}

		_ = t.c.Close()
		w.watchQueue.Remove(currentElement)
	}

	if w.watchQueue.Len() == 0 {
		w.scheduler.Stop()
		return
	}

	w.scheduler.Reset(time.Second)
}

func (w *ConnectionWatch) queueAdd(t *target) {
	w.watchQueue.PushBack(t)
	w.scheduler.Reset(time.Second)
}

func (w *ConnectionWatch) queueRemove(dst *target) {
	nextElement := w.watchQueue.Front()

	for nextElement != nil && nextElement.Value != nil {
		currentElement := nextElement
		nextElement = nextElement.Next()

		src := currentElement.Value.(*target)
		if src.id == dst.id {
			w.watchQueue.Remove(currentElement)
			break
		}
	}

	w.scheduler.Reset(time.Second)
}

// AddConnection .
func (w *ConnectionWatch) AddConnection(c net.Conn, maxBlockTime uint8) uint32 {
	id := atomic.AddUint32(&w.incrementID, 1)

	e := event{
		code: 1,
		content: &target{
			id:           id,
			c:            c,
			maxBlockTime: maxBlockTime,
		},
	}

	w.worker.Invoke(e, false)

	return id
}

// RemoveConnection .
func (w *ConnectionWatch) RemoveConnection(id uint32) {
	e := event{
		code:    2,
		content: &target{id: id},
	}

	w.worker.Invoke(e, false)
}
