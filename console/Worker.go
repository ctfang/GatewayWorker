package console

import (
	"GatewayWorker/events"
	"GatewayWorker/events/worker"
	"GatewayWorker/network"
	"GatewayWorker/network/tcp"
	"github.com/ctfang/command"
)

type Worker struct {
}

func (Worker) Configure() command.Configure {
	return command.Configure{
		Name:        "worker",
		Description: "业务worker进程",
		Input: command.Argument{
			Has: []command.ArgParam{
				{Name: "-d", Description: "是否使用守护进程"},
			},
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "register", Default: "127.0.0.1:1238", Description: "注册中心"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (Worker) Execute(input command.Input) {
	events.RegisterAddress = network.NewAddress(input.GetOption("register"))
	events.SecretKey = input.GetOption("secret")

	// 连接到注册中心
	register := tcp.NewClient()
	register.SetAddress(events.RegisterAddress)
	register.SetConnectionEvent(worker.NewRegisterEvent())
	register.ListenAndServe()
}
