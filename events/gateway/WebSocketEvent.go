package gateway

import (
	"GatewayWorker/events"
	"github.com/ctfang/network"
	"github.com/ctfang/network/protocol"
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
	// 内部通讯地址
	WorkerServerIp   string
	WorkerServerPort uint16
}

func (ws *WebSocketEvent) OnError(listen network.ListenTcp, err error) {

}

func (ws *WebSocketEvent) OnStart(listen network.ListenTcp) {
	ws.WorkerServerIp = events.WorkerAddress.Ip
	ws.WorkerServerPort = events.WorkerAddress.Port

	log.Println("ws server listening at: ", listen.GetAddress().Str)
}

func (ws *WebSocketEvent) GetClientId(client network.Connect) string {
	return network.Bin2hex(network.Ip2long(ws.WorkerServerIp), ws.WorkerServerPort, client.GetConnectionId())
}

func (ws *WebSocketEvent) OnConnect(client network.Connect) {
	client.SetUid(ws.GetClientId(client))
	// 添加连接池
	Router.AddedClient(client)

	header := &GatewayHeader{
		LocalIp:      network.Ip2long(ws.WorkerServerIp),
		LocalPort:    ws.WorkerServerPort,
		ClientIp:     client.GetIp(),
		ClientPort:   client.GetPort(),
		GatewayPort:  events.GatewayAddress.Port,
		ConnectionId: client.GetConnectionId(),
		flag:         1,
	}
	client.SetExtend(header)
	ws.SendToWorker(client, protocol.CMD_ON_CONNECT, []byte(""))
}

func (ws *WebSocketEvent) OnMessage(c network.Connect, message interface{}) {
	ws.SendToWorker(c, protocol.CMD_ON_MESSAGE, message.([]byte))
}

func (ws *WebSocketEvent) OnClose(c network.Connect) {
	ws.SendToWorker(c, protocol.CMD_ON_CLOSE, []byte(""))
	Router.DeleteClient(c.GetConnectionId())
}

func (ws *WebSocketEvent) SendToWorker(client network.Connect, cmd uint8, body []byte) {
	gh := client.GetExtend().(*GatewayHeader)
	msg := protocol.GatewayMessage{
		PackageLen:   28 + uint32(len(body)),
		Cmd:          cmd,
		LocalIp:      gh.LocalIp,
		LocalPort:    gh.LocalPort,
		ClientIp:     gh.ClientIp,
		ClientPort:   gh.ClientPort,
		ConnectionId: gh.ConnectionId,
		Flag:         gh.flag,
		GatewayPort:  gh.GatewayPort,
		ExtLen:       0,
		ExtData:      "",
		Body:         body,
	}

	worker, err := Router.GetWorker(client)
	if err != nil {
		// worker 找不到 获取连接
		log.Println("主动断开客户端连接 err:", err)
		client.Close()
		Router.DeleteClient(client.GetConnectionId())
		return
	}
	worker.Send(msg)
}

func NewWebSocketEvent() network.Event {
	return &WebSocketEvent{}
}
