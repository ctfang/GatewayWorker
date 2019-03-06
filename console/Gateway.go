package console

import (
	"GoGatewayWorker/gateway"
	"GoGatewayWorker/network"
	"github.com/ctfang/command"
)

type Gateway struct {
}

func (Gateway) Configure() command.Configure {
	return command.Configure{
		Name:        "gateway",
		Description: "网关进程gateway",
		Input: command.Argument{
			Has: []command.ArgParam{
				{Name: "-d", Description: "是否使用守护进程"},
			},
			Argument: []command.ArgParam{
				{Name: "runType", Description: "执行操作：start、stop、status"},
			},
			Option: []command.ArgParam{
				{Name: "gateway", Default: "127.0.0.1:8080", Description: "网关地址websocket"},
				{Name: "register", Default: "127.0.0.1:1238", Description: "注册中心"},
				{Name: "worker", Default: "127.0.0.1:4000", Description: "内部通讯地址"},
				{Name: "secret", Default: "", Description: "通讯秘钥"},
			},
		},
	}
}

func (Gateway) Execute(input command.Input) {

	gateway.GatewayAddress = network.NewAddress(input.GetOption("gateway"))
	gateway.WorkerAddress = network.NewAddress(input.GetOption("worker"))
	gateway.RegisterAddress = network.NewAddress(input.GetOption("register"))
	gateway.SecretKey = input.GetOption("secret")

	ws := network.WebSocket{
		Addr:  gateway.GatewayAddress,
		Event: &gateway.GatewayEvent{},
	}
	ws.ListenAndServe()
}
