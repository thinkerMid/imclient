package imclient

import (
	"context"
	"github.com/hyahm/golog"
	"labs/src/imapp"
	"labs/src/imapp/wsa"
)

type IMClient struct {
	im imapp.IM

	cancel      context.CancelFunc
	commandChan chan CmdId
	eventChan   chan EventId
}

func CreateClient() *IMClient {
	command := make(chan CmdId, 1)
	event := make(chan EventId, 1024)

	return &IMClient{
		im:          wsa.CreateWSApp(event),
		cancel:      nil,
		commandChan: command,
		eventChan:   event,
	}
}

func (cli *IMClient) RunClient(country string) {
	defer func() {
		golog.Infof("client run over")
	}()
	golog.Infof("client run start")

	go cli.im.Launch()

	// 处理命令
	func() {
		ctx, cancel := context.WithCancel(context.Background())
		cli.cancel = cancel

		for true {
			select {
			case <-ctx.Done():
				golog.Infof("IM:%v run over\n", cli.im.Info())
				return
			case cmdId := <-cli.commandChan:
				cli.processCommand(cmdId)
			}
		}
	}()
}

func (cli *IMClient) Cancel() {
	if cli.im != nil {
		cli.im.Kill()
	}
	if cli.cancel != nil {
		cli.cancel()
	}
}

func (cli *IMClient) processCommand(cmdId CmdId) {
	golog.Infof("user:%v process command:%v\n", cli.im.Info(), cmdId)
}
