package lightning

import (
	"context"
	"github.com/marusama/semaphore/v2"
	"sync/atomic"
	"time"
)

const singleWorker int32 = 1

// 等待时间 / 毫秒
const waitTime = time.Millisecond * 10

type rwState uint8

const (
	emptyState rwState = iota
	writeState
)

// WorkPool .
type WorkPool struct {
	ctx    context.Context
	cancel context.CancelFunc

	// 协程资源锁
	semaphoreG semaphore.Semaphore

	workerCount int32

	taskLimitIdx int32 // 任务储存数
	readIdx      int32 // 读下标
	writeIdx     int32 // 写下标

	queue      []interface{} // 任务储存队列
	stateQueue []rwState     //  任务储存队列所对应的可读写状态

	advanceDeliveryCount int32 // 预投递任务数	*用于表示当前任务池的未来负载数量
	deliveredTaskCount   int32 // 已投递的任务数	*当前任务池已承载的数量

	readGoroutineStatus int32 // 读协程运行状态

	workerFn func(interface{})
}

// New .
func New(worker int32, taskLimit int32, workerFn func(interface{})) *WorkPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkPool{
		ctx:          ctx,
		cancel:       cancel,
		workerCount:  worker,
		taskLimitIdx: taskLimit - 1,
		queue:        make([]interface{}, taskLimit),
		stateQueue:   make([]rwState, taskLimit),
		workerFn:     workerFn,
		semaphoreG:   semaphore.New(int(worker)),
	}
}

// NewAsyncPool 异步事件循环的池子
func NewAsyncPool(taskLimit int32, workerFn func(interface{})) *WorkPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkPool{
		ctx:          ctx,
		cancel:       cancel,
		workerCount:  singleWorker,
		taskLimitIdx: taskLimit - 1,
		queue:        make([]interface{}, taskLimit),
		stateQueue:   make([]rwState, taskLimit),
		workerFn:     workerFn,
		semaphoreG:   semaphore.New(int(singleWorker)),
	}
}

// Stop .
func (r *WorkPool) Stop() {
	r.cancel()

	_ = r.semaphoreG.Acquire(context.Background(), r.semaphoreG.GetLimit())
	r.semaphoreG.Release(r.semaphoreG.GetLimit())
}

// Await .
func (r *WorkPool) Await() {
	for {
		if atomic.LoadInt32(&r.deliveredTaskCount) == 0 {
			break
		}

		time.Sleep(waitTime)
	}

	_ = r.semaphoreG.Acquire(context.Background(), r.semaphoreG.GetLimit())
	r.semaphoreG.Release(r.semaphoreG.GetLimit())
}

// Invoke .
//
// 任务池假定为的满的情况 都使用go原语执行 避免消费协程进入投递操作 导致生产和消费协程死循环
// 以此使用更多的协程来代替执行协程投递任务 确保当前执行的协程执行完成 正常消费任务
//
// *上述的情况 不可能出现于正常的生产消费模式 即为消费者是不会变成生产者进行任务投递
//
//	以下是为了解决 在单协程异步事件循环的并发场景下 消费协程会进行任务队列尾投递操作 如果任务队列满了会进入无限等待投递的过程 这样子是死循环
func (r *WorkPool) Invoke(i interface{}, async bool) {
	// 预投递
	advanceDeliveryCount := atomic.AddInt32(&r.advanceDeliveryCount, 1)

	// 异步
	if async {
		// 下一次投递会占满任务
		if advanceDeliveryCount >= r.taskLimitIdx {
			go r.delivery(i)
			return
		}
	}

	r.delivery(i)
}

func (r *WorkPool) delivery(i interface{}) {
	var newWriteIdx int32
	var nowWriteIdx int32

	for {
		if r.done() {
			atomic.AddInt32(&r.advanceDeliveryCount, -1)
			break
		}

		nowWriteIdx = atomic.LoadInt32(&r.writeIdx)

		if r.stateQueue[nowWriteIdx] == emptyState {
			if nowWriteIdx == r.taskLimitIdx {
				newWriteIdx = 0
			} else {
				newWriteIdx = nowWriteIdx + 1
			}

			if !atomic.CompareAndSwapInt32(&r.writeIdx, nowWriteIdx, newWriteIdx) {
				continue
			}

			r.queue[nowWriteIdx] = i
			r.stateQueue[nowWriteIdx] = writeState
		} else {
			// 使用 runtime.Gosched() 的初衷
			//   这里使用了较为激进的方案，主是为了保持高性能低延时，比sleep会更加占用CPU
			//
			// 如果这里出现了高频调用造成CPU占用过高，那么需要排查两个地方，这些都是造成高频调用的原因
			//  上游：
			// 			1.投递速度过快，下游消费能力不足
			//  下游：
			// 			1.消费协程处理时间过长，导致上游投递堆积并且造成高频调用等待投递
			//			2.协程数量不够，导致上游投递堆积并且造成高频调用等待投递
			//
			//  统一的解决方案：进行参数调优，调整协程数量和任务容量（协程数量的增加意味着CPU占用会增加，任务容量的增加意味着内存占用会增加）
			//           调整到一个合适的参数可以缓解高频调用的问题
			//  runtime.Gosched()

			// 2022-10-3 等待投递的实现方式变更
			//  据线上数据处理量 CPU使用率等现象表明 只会出现个别协程处理的任务量爆满 容易导致就这几个协程来回空转 cpu使用率比较高
			//  改用睡眠操作后，会在协程数量多的时候会导致协程唤醒的时间是存在不确定性的，runtime调度足够快的时候不存在这种现象
			//  和 runtime.Gosched() 相比带来的好处就是：cpu使用率会比之前更低。但是存在任务处理效率有稍微下降的可能性，因为睡眠的时间是不准确的
			time.Sleep(waitTime)
			continue
		}

		// 投递成功的任务数增加
		atomic.AddInt32(&r.deliveredTaskCount, 1)

		if atomic.CompareAndSwapInt32(&r.readGoroutineStatus, 0, 1) {
			go r.processQueue()
		}

		break
	}
}

func (r *WorkPool) processQueue() {
	var newReadIdx int32
	var nowReadIdx int32

	for {
		nowReadIdx = atomic.LoadInt32(&r.readIdx)

		if r.stateQueue[nowReadIdx] == writeState {
			if nowReadIdx == r.taskLimitIdx {
				newReadIdx = 0
			} else {
				newReadIdx = nowReadIdx + 1
			}

			if !atomic.CompareAndSwapInt32(&r.readIdx, nowReadIdx, newReadIdx) {
				continue
			}

			if r.workerCount == singleWorker {
				r.workerFn(r.queue[nowReadIdx])
			} else {
				_ = r.semaphoreG.Acquire(context.Background(), 1)

				go r.runGoroutine(r.queue[nowReadIdx])
			}

			// 清空
			r.queue[nowReadIdx] = emptyState
			r.stateQueue[nowReadIdx] = emptyState

			// 减去预投递数量
			atomic.AddInt32(&r.advanceDeliveryCount, -1)

			// 减去已投递数量
			if atomic.AddInt32(&r.deliveredTaskCount, -1) > 0 {
				continue
			}
		}

		atomic.StoreInt32(&r.readGoroutineStatus, 0)

		if atomic.LoadInt32(&r.deliveredTaskCount) > 0 && atomic.CompareAndSwapInt32(&r.readGoroutineStatus, 0, 1) {
			continue
		}

		break
	}
}

func (r *WorkPool) runGoroutine(i interface{}) {
	r.workerFn(i)
	r.semaphoreG.Release(1)
}

func (r *WorkPool) done() bool {
	select {
	case <-r.ctx.Done():
		return true
	default:
		return false
	}
}
