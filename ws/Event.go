package ws

import (
	"GatewayWorker/events/worker"
)

type Event struct {
}

func NewEvent() *Event {
	return &Event{}
}

func (*Event) OnStart() {
	panic("implement me")
}

func (*Event) OnConnect(client worker.Connect) {
	panic("implement me")
}

func (*Event) OnMessage(client worker.Connect, message interface{}) {
	panic("implement me")
}

func (*Event) OnClose(client worker.Connect) {
	panic("implement me")
}
