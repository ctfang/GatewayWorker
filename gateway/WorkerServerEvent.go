package gateway

import (
	"GoGatewayWorker/network"
	"GoGatewayWorker/protocol"
	"fmt"
	"log"
	"time"
)

/*
worker内部通讯逻辑,接收worker发到客户端信息
 */
type WorkerServerEvent struct {
	SendToWorker chan []byte
}

func (w *WorkerServerEvent) OnStart(tcp *network.TcpServer) {
	go w.connectRegister()
	w.SendToWorker = make(chan []byte,256)
	log.Println("已启动内部通讯",WorkerAddress.Str)
}

/*
worker已连接网关
 */
func (*WorkerServerEvent) OnConnect(c *network.TcpServerClient) {

}

func (ws *WorkerServerEvent) OnMessage(c *network.TcpServerClient, message interface{}) {
	msg := message.(protocol.GatewayMessage)

	switch msg.Cmd {
	case protocol.CMD_ON_CONNECT:

	case protocol.CMD_ON_MESSAGE:
	case protocol.CMD_ON_CLOSE:
		ws.OnClose(c)
	case protocol.CMD_SEND_TO_ONE:
		// 单个用户信息
		ws.sendToOne(msg)
	case protocol.CMD_SEND_TO_ALL:
		// 发给gateway的向所有用户发送数据
		ws.sendToAll(msg)
	case protocol.CMD_WORKER_CONNECT:
		Router.AddedWorker(c)
		log.Println("worker已连接网关")
	default:
		log.Println(message)
	}
}

func (*WorkerServerEvent) sendToOne(msg protocol.GatewayMessage) {
	client,err := Router.GetClient(msg.ConnectionId)
	if err!=nil {
		Router.DeleteClient(msg.ConnectionId)
		return
	}

	client.Send([]byte(msg.Body))
}

func (*WorkerServerEvent) sendToAll(msg protocol.GatewayMessage) {
	for _,client := range Router.Clients {
		client.Send([]byte(msg.Body))
	}
}

/*
关闭
 */
func (w *WorkerServerEvent) OnClose(c *network.TcpServerClient) {
	Router.DeleteWorker(c)
	_ = c.Close()
}


func (w *WorkerServerEvent) SendToWorkerTask() {
	log.Println("implement me")
}

func (w *WorkerServerEvent) sendToWorker() {
	for {
		select {
		case message := <- w.SendToWorker:
			fmt.Println(message)
		}
	}
}

// 连接到注册中心
func (w *WorkerServerEvent) connectRegister() {
	tcp := &network.TcpServerConnection{
		Addr:     RegisterAddress,
		Event:    &RegisterEvent{},
		Protocol: protocol.Text{},
		SendChan: make(chan []byte, 256),
	}

	err := tcp.Connect()
	if err != nil {
		// 连接失败，添加定时器，定时请求
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ticker.C:
				err := tcp.Connect()
				if err == nil {
					ticker.Stop()
					return
				}
			}
		}
	}
}
