package watchDog

import (
	"container/list"
	"fmt"
	"go.uber.org/zap"
	"ws/framework/application/container/abstract_interface"
	"ws/framework/application/core/processor"
)

const (
	noneLevel uint8 = iota
	blockInfoLevel
	blockWarnLevel
	blockErrorLevel
)

const (
	start uint8 = iota
	receive
)

var startTemplate = "the `%s` action has been blocked in [Start] status for at least %v seconds"
var receiveTemplate = "the `%s` action has been blocked in [Receive-%s] status for at least %v seconds"
var ignoreAction = map[string]struct{}{
	"Wait":                   {},
	"Offline":                {},
	"SimulateInputChatState": {},
}

// ----------------------------------------------------------------------------

type checkResult struct {
	actionIndex       uint8  // 上次的action
	actionName        string // 结构体的名字
	lastExecuteStatus uint8  // 执行状态
	blockTime         uint8  // 阻塞时间 单位秒
}

// WatchDog .
type WatchDog struct {
	Logger *zap.SugaredLogger

	processorCheckResultMapping map[uint32]checkResult

	processorLinkList *list.List
}

// New .
func New(logger *zap.SugaredLogger) *WatchDog {
	return &WatchDog{
		Logger:                      logger.Named("WatchDog"),
		processorCheckResultMapping: make(map[uint32]checkResult),
		processorLinkList:           list.New(),
	}
}

// Watch .
func (w *WatchDog) Watch(processor containerInterface.IProcessor) {
	w.processorLinkList.PushBack(processor)
}

// UnWatch .
func (w *WatchDog) UnWatch(processor containerInterface.IProcessor) {
	removePID := processor.ID()

	element := w.processorLinkList.Front()

	for element != nil && element.Value != nil {
		p := element.Value.(containerInterface.IProcessor)
		if removePID == p.ID() {
			w.processorLinkList.Remove(element)
			delete(w.processorCheckResultMapping, p.ID())
			break
		}

		element = element.Next()
	}
}

// ScheduleCheck .
func (w *WatchDog) ScheduleCheck() error {
	element := w.processorLinkList.Front()

	var levelCode uint8
	var err error

	for element != nil && element.Value != nil {
		p := element.Value.(containerInterface.IProcessor)
		if p.ProcessorType() == containerInterface.LogicProcessorType {
			levelCode, err = w.checkLogicProcessor(p.(*processor.LogicProcessor))

			switch levelCode {
			case blockInfoLevel:
				w.Logger.Info(err)
			case blockWarnLevel:
				w.Logger.Warn(err)
			case blockErrorLevel:
				// 这个等级了才抛出
				w.Logger.Error(err)
				return err
			}
		}

		element = element.Next()
	}

	return nil
}

func (w *WatchDog) checkLogicProcessor(p *processor.LogicProcessor) (levelCode uint8, err error) {
	/*
		1.检测action处于什么状态，启动和接收耗时是否超时
		2.前台和后台模式 前台30秒 后台10秒
	*/

	// 触发器未激活的 不做监控
	if p.Trigger.WaitActive() {
		return
	}

	res := w.processorCheckResultMapping[p.ID()]
	action := p.ActionQueue.Current()

	var nowStatus uint8
	if len(action.ReceiveID()) > 0 {
		nowStatus = receive
	}

	// 和上一次检测的action同个ID
	if res.actionIndex == action.ActionIndex() {
		if nowStatus == res.lastExecuteStatus {
			res.blockTime = res.blockTime + 1
		} else {
			res.blockTime = 0
		}
	} else {
		res.actionIndex = action.ActionIndex()
		res.lastExecuteStatus = nowStatus
		res.blockTime = 0
		res.actionName = action.ActionName()
	}

	w.processorCheckResultMapping[p.ID()] = res

	// 属于特殊的Action
	if _, ok := ignoreAction[res.actionName]; ok {
		levelCode = noneLevel
		return
		//// 减少对这种的action报警 每60秒才警示一次 其余忽略
		//if res.blockTime%60 > 0 {
		//	levelCode = noneLevel
		//	return
		//}
		//
		//// 只显示INFO级别
		//levelCode = blockInfoLevel
	}

	levelCode = w.computeLevel(res.blockTime, p.Priority() == processor.PriorityBackground)
	if levelCode == noneLevel {
		return
	}

	if res.lastExecuteStatus == start {
		err = fmt.Errorf(startTemplate, res.actionName, res.blockTime)
	} else {
		err = fmt.Errorf(receiveTemplate, res.actionName, action.ReceiveID(), res.blockTime)
	}

	return
}

func (w *WatchDog) computeLevel(blockTime uint8, isBackground bool) (levelCode uint8) {
	// 后台的 10秒info 15秒warn 25秒error
	if isBackground {
		if blockTime == 25 {
			levelCode = blockErrorLevel
		} else if blockTime >= 15 {
			levelCode = blockWarnLevel
		} else if blockTime == 10 {
			levelCode = blockInfoLevel
		}

		return
	}

	// 前台的 10秒info 25秒warn 35秒error
	if blockTime == 35 {
		levelCode = blockErrorLevel
	} else if blockTime >= 25 {
		levelCode = blockWarnLevel
	} else if blockTime == 10 {
		levelCode = blockInfoLevel
	}

	return
}

// Stop .
func (w *WatchDog) Stop() {
	w.processorLinkList.Init()
}
