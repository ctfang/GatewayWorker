package worker

import (
	"GatewayWorker/events"
	"GatewayWorker/network"
	"GatewayWorker/network/protocol"
)

type GatewayHandle struct {
}

func (g *GatewayHandle) OnConnect(message protocol.GatewayMessage) {
	clientId := network.Bin2hex(message.LocalIp, message.LocalPort, message.ConnectionId)
	events.BussinessEvent.OnConnect(clientId)
}

func (*GatewayHandle) OnMessage(message protocol.GatewayMessage) {
	clientId := network.Bin2hex(message.LocalIp, message.LocalPort, message.ConnectionId)
	events.BussinessEvent.OnMessage(clientId, message.Body)
}

func (*GatewayHandle) OnClose(message protocol.GatewayMessage) {
	clientId := network.Bin2hex(message.LocalIp, message.LocalPort, message.ConnectionId)
	events.BussinessEvent.OnClose(clientId)
}
