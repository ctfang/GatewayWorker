package gateway

import (
	"GoGatewayWorker/network"
	"GoGatewayWorker/protocol"
	"encoding/binary"
	"encoding/hex"
	"log"
)

/*
网关逻辑
转发数据到worker
 */
type GatewayEvent struct {
	WorkerServer *WorkerServerEvent
}

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

func (g *GatewayEvent) OnStart() {
	// 启动一个内部通讯tcp server
	tcp := network.TcpServer{}
	worker := &WorkerServerEvent{}
	g.WorkerServer = worker
	tcp.SetAddress(WorkerAddress)
	tcp.SetProtocol(&protocol.GatewayProtocol{})
	tcp.SetEvent(worker)
	tcp.ListenAndServe()
}

/*
有客户端连接
 */
func (g *GatewayEvent) OnConnect(client *network.TcpClientConnection) {
	client.SetUid(bin2hex(WorkerAddress.Ip,WorkerAddress.Port,client.GetConnectionId()))
	_, err := Router.AddedClient(client)
	if err != nil {
		log.Fatalln(err)
		g.OnClose(client)
		return
	}
	header := GatewayHeader{
		LocalIp:      network.Ip2long(WorkerAddress.Ip),
		LocalPort:    WorkerAddress.Port,
		ClientIp:     network.Ip2long(client.GetIp()),
		ClientPort:   client.GetPort(),
		GatewayPort:  GatewayAddress.Port,
		ConnectionId: client.GetConnectionId(),
		flag:         1,
	}
	client.Extend = header
	g.SendToWorker(client, protocol.CMD_ON_CONNECT, "")
}

// 构建分布式唯一id
func bin2hex(ip string, port uint16, id uint32)string{
	var msgByte []byte
	ipUint32 := network.Ip2long(ip)
	var buf32 = make([]byte, 4)
	var bug16 = make([]byte, 2)
	binary.BigEndian.PutUint32(buf32, ipUint32)
	msgByte = append(msgByte, buf32...)
	binary.BigEndian.PutUint16(bug16, port)
	msgByte = append(msgByte, bug16...)
	binary.BigEndian.PutUint32(buf32, id)
	msgByte = append(msgByte, buf32...)
	return hex.EncodeToString(msgByte)
}

// 客户端信息转发到worker处理
func (g *GatewayEvent) OnMessage(clint *network.TcpClientConnection, message []byte) {
	body := string(message)
	g.SendToWorker(clint, protocol.CMD_ON_MESSAGE, body)
}

func (GatewayEvent) OnClose(clint *network.TcpClientConnection) {
	Router.DeleteClient(clint.GetConnectionId())
	clint.Close()
}

func (g *GatewayEvent) SendToWorker(client *network.TcpClientConnection, cmd uint8, body string) {
	GatewayHeader := client.Extend.(GatewayHeader)
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