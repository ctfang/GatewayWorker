package console

import (
	"GatewayWorker/network"
	"GatewayWorker/network/tcp"
	"github.com/ctfang/command"
	"log"
)

type Test struct {
}

func (Test) Configure() command.Configure {
	return command.Configure{
		Name:        "test",
		Description: "test",
		Input:       command.Argument{},
	}
}

func (Test) Execute(input command.Input) {
	server := tcp.Server{}
	server.SetAddress(network.NewAddress(":8080/ws"))
	server.SetConnectionEvent(&Hello{})
	server.ListenAndServe()
}

type Hello struct {
}

func (*Hello) OnError(listen network.ListenTcp, err error) {
	panic("implement me")
}

func (*Hello) OnStart(tcp network.ListenTcp) {
	log.Println("OnStart")
}

func (*Hello) OnConnect(c network.Connect) {
	log.Println("OnConnect")
}

func (*Hello) OnMessage(c network.Connect, message interface{}) {
	c.Send([]byte("收到信息"))
	c.Send([]byte("1"))
	c.Send([]byte("b"))
	c.Send([]byte("3"))
	log.Println(message)
}

func (*Hello) OnClose(c network.Connect) {
	log.Println("OnClose")
}
