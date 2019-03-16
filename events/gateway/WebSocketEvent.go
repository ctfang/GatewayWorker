package gateway

import (
	"GatewayWorker/events"
	"GatewayWorker/network"
	"GatewayWorker/network/protocol"
	"log"
)

/**
客户端信息
*/
type GatewayHeader struct {
	// 内部通讯地址 , 对应本机地址
	LocalIp      uint32
	LocalPort    uint16
	ClientIp     uint32
	ClientPort   uint16
	GatewayPort  uint16
	ConnectionId uint32
	flag         uint8
}

type WebSocketEvent struct {
}

func (*WebSocketEvent) OnError(listen network.ListenTcp, err error) {

}

func (ws *WebSocketEvent) OnStart(listen network.ListenTcp) {
	log.Println("ws server listening at: ", listen.GetAddress().Str)
}

func (ws *WebSocketEvent) OnConnect(client network.Connect) {
	client.SetUid(network.Bin2hex(events.WorkerAddress.Ip, events.WorkerAddress.Port, client.GetConnectionId()))
	_, err := Router.AddedClient(client)
	if err != nil {
		log.Println(err)
		ws.OnClose(client)
		return
	}
	header := &GatewayHeader{
		LocalIp:      network.Ip2long(events.WorkerAddress.Ip),
		LocalPort:    events.WorkerAddress.Port,
		ClientIp:     network.Ip2long(client.GetIp()),
		ClientPort:   client.GetPort(),
		GatewayPort:  events.GatewayAddress.Port,
		ConnectionId: client.GetConnectionId(),
		flag:         1,
	}
	client.SetExtend(header)
	ws.SendToWorker(client, protocol.CMD_ON_CONNECT, "")
}

func (ws *WebSocketEvent) OnMessage(c network.Connect, message interface{}) {
	body := string(message.([]byte))
	ws.SendToWorker(c, protocol.CMD_ON_MESSAGE, body)
}

func (*WebSocketEvent) OnClose(c network.Connect) {
	Router.DeleteClient(c.GetConnectionId())
	c.Close()
}

func (ws *WebSocketEvent) SendToWorker(client network.Connect, cmd uint8, body string) {
	GatewayHeader := client.GetExtend().(GatewayHeader)
	msg := protocol.GatewayMessage{
		PackageLen:   28 + uint32(len(body)),
		Cmd:          cmd,
		LocalIp:      GatewayHeader.LocalIp,
		LocalPort:    GatewayHeader.LocalPort,
		ClientIp:     GatewayHeader.ClientIp,
		ClientPort:   GatewayHeader.ClientPort,
		ConnectionId: GatewayHeader.ConnectionId,
		Flag:         GatewayHeader.flag,
		GatewayPort:  GatewayHeader.GatewayPort,
		ExtLen:       0,
		ExtData:      "",
		Body:         body,
	}

	worker := Router.GetWorker(client)
	worker.Send(msg)
}

func NewWebSocketEvent() network.Event {
	return &WebSocketEvent{}
}
