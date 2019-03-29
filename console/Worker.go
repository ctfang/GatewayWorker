package console

import (
	"GatewayWorker/events"
	"GatewayWorker/events/worker"
	"GatewayWorker/network"
	"GatewayWorker/network/tcp"
	"github.com/ctfang/command"
	"log"
	"time"
)

type Worker struct {
}

func (Worker) Configure() command.Configure {
	return command.Configure{
		Name:        "worker",
		Description: "业务worker进程",
		Input: command.Argument{
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
	events.BussinessEvent = NewHelloEvent()

	// 连接到注册中心
	register := tcp.NewClient()
	register.SetAddress(events.RegisterAddress)
	register.SetConnectionEvent(worker.NewRegisterEvent())
	register.ListenAndServe()

	// 断线重连
	for {
		ticker := time.NewTicker(time.Second * 2)
		select {
		case <-ticker.C:
			register.ListenAndServe()
		}
	}

}

type HelloEvent struct {
}

func NewHelloEvent() *HelloEvent {
	event := HelloEvent{}
	event.OnStart()
	return &event
}

func (*HelloEvent) OnStart() {
	log.Println("OnStart")
}

func (*HelloEvent) OnConnect(clientId string) {
	log.Println("OnConnect ", clientId)
}

func (*HelloEvent) OnMessage(clientId string, body []byte) {
	log.Println("OnMessage ", clientId, string(body))
}

func (*HelloEvent) OnClose(clientId string) {
	log.Println("OnClose ", clientId)
}
