package console

import (
	"GatewayWorker/events"
	"GatewayWorker/events/register"
	"GatewayWorker/network"
	"GatewayWorker/network/tcp"
	"github.com/ctfang/command"
)

type Register struct {
}

func (Register) Configure() command.Configure {
	return command.Configure{
		Name:        "register",
		Description: "注册中心",
		Input: command.Argument{
			Has: []command.ArgParam{
				{Name: "-d", Description: "是否使用守护进程"},
			},
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "register", Description: "手动设置地址"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (Register) Execute(input command.Input) {
	events.RegisterAddress = network.NewAddress(input.GetOption("register"))
	events.SecretKey = input.GetOption("secret")

	server := tcp.NewServer()
	server.SetAddress(events.RegisterAddress)
	server.SetConnectionEvent(register.NewRegisterEvent())
	server.ListenAndServe()
}
