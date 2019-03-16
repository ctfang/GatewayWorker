package console

import (
	"GatewayWorker/events"
	"GatewayWorker/events/gateway"
	"GatewayWorker/network"
	"GatewayWorker/network/protocol"
	"GatewayWorker/network/tcp"
	"GatewayWorker/network/ws"
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
	events.GatewayAddress = network.NewAddress(input.GetOption("gateway"))
	events.WorkerAddress = network.NewAddress(input.GetOption("worker"))
	events.RegisterAddress = network.NewAddress(input.GetOption("register"))
	events.SecretKey = input.GetOption("secret")

	// 启动一个内部通讯tcp server
	worker := tcp.NewServer()
	worker.SetAddress(events.WorkerAddress)
	worker.SetConnectionEvent(gateway.NewWorkerEvent())
	worker.SetProtocol(protocol.NewGatewayProtocol())
	go worker.ListenAndServe()

	// 连接到注册中心
	register := tcp.NewClient()
	register.SetAddress(events.RegisterAddress)
	register.SetConnectionEvent(gateway.NewRegisterEvent())
	go register.ListenAndServe()

	// 启动对客户端的websocket连接
	server := ws.Server{}
	server.SetAddress(events.GatewayAddress)
	server.SetConnectionEvent(gateway.NewWebSocketEvent())
	server.ListenAndServe()
}
