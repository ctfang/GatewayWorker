package console

import (
	"GoGatewayWorker/gateway"
	"GoGatewayWorker/network"
	"GoGatewayWorker/protocol"
	"GoGatewayWorker/register"
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
				{Name: "addr", Description: "手动设置地址"},
			},
		},
	}
}

func (Register) Execute(input command.Input) {
	gateway.RegisterAddress = network.NewAddress(input.GetOption("register"))
	gateway.SecretKey = input.GetOption("secret")

	tcp := network.TcpServer{}
	tcp.SetAddress(gateway.RegisterAddress)
	tcp.SetProtocol(protocol.Text{})
	tcp.SetEvent(register.NewRegisterServerEvent())
	tcp.ListenAndServe()
}
