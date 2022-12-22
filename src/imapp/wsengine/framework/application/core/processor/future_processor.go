package processor

import (
	"go.uber.org/zap"
	"ws/framework/application/container/abstract_interface"
	functionTools "ws/framework/utils/function_tools"
)

// FutureProcessor .
type FutureProcessor struct {
	BaseProcessor

	logger       *zap.SugaredLogger
	waitFutureID string
	future       containerInterface.IFuture
}

// NewFutureProcessor .
func NewFutureProcessor(waitFutureID string, future containerInterface.IFuture) containerInterface.IProcessor {
	return &FutureProcessor{
		waitFutureID: waitFutureID,
		future:       future,
	}
}

// Init .
func (f *FutureProcessor) Init(logger *zap.SugaredLogger) {
	f.logger = logger
}

// ProcessMessage .
func (f *FutureProcessor) ProcessMessage(context containerInterface.IProcessControl) {
	if len(f.waitFutureID) == 0 {
		context.Reject()
		return
	}

	if f.waitFutureID != context.MessageContext().Message().ID() {
		return
	}

	f.actionReceiveNode(context)
}

func (f *FutureProcessor) actionReceiveNode(context containerInterface.IProcessControl) {
	msgCtx := context.MessageContext()
	err := msgCtx.Message().HasError()

	if err == nil {
		err = f.future.Receive(msgCtx)
	}

	if err != nil {
		f.future.Error(msgCtx, err)
		f.logger.Errorf("%s `%s` future wait result (%s) error", "FutureProcessorType", functionTools.ReflectValueTypeName(f.future), err)
	}

	context.Reject()
}

// OnChannelError .
func (f *FutureProcessor) OnChannelError(context containerInterface.IProcessControl, err error) {
	f.future.Error(context.MessageContext(), err)
	f.logger.Errorf("%s `%s` future wait result (%s) error", "FutureProcessorType", functionTools.ReflectValueTypeName(f.future), err)
}

// ProcessorType .
func (f *FutureProcessor) ProcessorType() containerInterface.ProcessorType {
	return containerInterface.FutureProcessorType
}

// DumpInfo .
func (f *FutureProcessor) DumpInfo() containerInterface.ProcessorDump {
	dump := containerInterface.ProcessorDump{
		PID:           f.pid,
		AliasName:     f.aliasName,
		ProcessorType: f.ProcessorType().String(),
		ActionQueue:   "/",
	}

	dump.CurrentActionName = functionTools.ReflectValueTypeName(f.future)
	dump.CurrentActionStatus = "Receive"

	return dump
}
