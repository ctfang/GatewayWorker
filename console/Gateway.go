package console

import (
	"GatewayWorker/events"
	"GatewayWorker/events/gateway"
	"fmt"
	"github.com/ctfang/command"
	"github.com/ctfang/network"
	"github.com/ctfang/network/protocol"
	"github.com/ctfang/network/tcp"
	"github.com/ctfang/network/ws"
	"github.com/gorilla/websocket"
	"log"
)

type Gateway struct {
	Name string
}

func (self *Gateway) Configure() command.Configure {
	self.Name = "gateway"
	return command.Configure{
		Name:        self.Name,
		Description: "网关进程gateway",
		Input: command.Argument{
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

func (self *Gateway) Execute(input command.Input) {
	switch input.GetArgument("runType") {
	case "start":
		self.start(input)
	case "stop":
		self.stop(input)
	case "status":
		self.status(input)
	}
}

func (self *Gateway) start(input command.Input) {
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
	ws.MessageType = websocket.BinaryMessage
	server := ws.Server{}
	server.SetAddress(events.GatewayAddress)
	server.SetConnectionEvent(gateway.NewWebSocketEvent())
	server.ListenAndServe()
}

func (self *Gateway) status(input command.Input) {
	log.Println("未做")
}

func (self *Gateway) stop(input command.Input) {
	err := command.StopSignal(self.Name)
	if err != nil {
		fmt.Println("停止失败")
	}
	fmt.Println("停止成功")
}
