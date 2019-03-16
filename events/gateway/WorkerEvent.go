package gateway

import (
	"GatewayWorker/network"
	"log"
)

type WorkerEvent struct {
}

func (*WorkerEvent) OnError(listen network.ListenTcp, err error) {

}

func (*WorkerEvent) OnStart(listen network.ListenTcp) {
	log.Println("worker server listening at: ", listen.GetAddress().Str)
}

func (*WorkerEvent) OnConnect(c network.Connect) {

}

func (*WorkerEvent) OnMessage(c network.Connect, message interface{}) {
	panic("implement me")
}

func (*WorkerEvent) OnClose(c network.Connect) {
	panic("implement me")
}

func NewWorkerEvent() network.Event {
	return &WorkerEvent{}
}
