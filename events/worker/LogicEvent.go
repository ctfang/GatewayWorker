package worker

import (
	"GatewayWorker/events"
	"GatewayWorker/network"
	"log"
)

type WorkerKey struct {
	Worker string `json:"worker_key"`
	Secret string `json:"secret_key"`
}

type LogicEvent struct {
	registerEvent    *RegisterEvent
	gatewayAddresses string
}

func (*LogicEvent) OnStart(listen network.ListenTcp) {

}

func (l *LogicEvent) OnConnect(c network.Connect) {
	l.registerEvent.UpdateGatewayConnections(l.gatewayAddresses, c)
	msg := WorkerKey{
		Worker: "BussinessWorker:" + string(c.GetConnectionId()),
		Secret: events.SecretKey,
	}
	c.Send(msg)
}

func (*LogicEvent) OnMessage(c network.Connect, message interface{}) {
	log.Println(message)
}

func (l *LogicEvent) OnClose(c network.Connect) {
	l.registerEvent.UpdateGatewayConnections(l.gatewayAddresses, nil)
}

func (l *LogicEvent) OnError(listen network.ListenTcp, err error) {
	l.registerEvent.UpdateGatewayConnections(l.gatewayAddresses, nil)
}

func NewLogicEvent(r *RegisterEvent, gatewayAddresses string) network.Event {
	return &LogicEvent{
		registerEvent:    r,
		gatewayAddresses: gatewayAddresses,
	}
}
