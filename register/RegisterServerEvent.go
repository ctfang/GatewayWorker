package register

import (
	"GoGatewayWorker/network"
	"encoding/json"
	"fmt"
)

// 注册中心
type RegisterServerEvent struct {
	// ID 地址保存
	gatewayConnections map[uint32]string
	workerConnections  map[uint32]*network.TcpServerClient
	// 秘要
	SecretKey string
}

func NewRegisterServerEvent() *RegisterServerEvent {
	reg := RegisterServerEvent{
		gatewayConnections: make(map[uint32]string),
		workerConnections:  make(map[uint32]*network.TcpServerClient),
	}
	return &reg
}

type RegisterMessage struct {
	Event     string `json:"event"`
	Address   string `json:"address"`
	SecretKey string `json:"secret_key"`
}


func (r *RegisterServerEvent) OnStart(tcp *network.TcpServer) {

}

// 新链接
func (r *RegisterServerEvent) OnConnect(c *network.TcpServerClient) {
	//
	c.Send("1")
	c.Send("2")
}

// 新信息
func (r *RegisterServerEvent) OnMessage(c *network.TcpServerClient, msg interface{}) {
	var data RegisterMessage
	err := json.Unmarshal([]byte(msg.(string)), &data)
	if err != nil {
		fmt.Println(err)
		_ = c.Close()
		return
	}
	if r.SecretKey != "" {
		if data.SecretKey != r.SecretKey {
			fmt.Println("秘要不对")
			_ = c.Close()
			return
		}
	}

	switch data.Event {
	case "gateway_connect":
		r.gatewayConnect(c, data)
	case "worker_connect":
		r.workerConnect(c, data)
	case "ping":
		return
	default:
		fmt.Println("不认识的事件定义")
		_ = c.Close()
	}
}

// 链接关闭
func (r *RegisterServerEvent) OnClose(c *network.TcpServerClient) {
	_, hasG := r.gatewayConnections[c.Id]
	if hasG == true {
		delete(r.gatewayConnections, c.Id)
		r.broadcastAddresses(0)
	}

	_, hasW := r.workerConnections[c.Id]
	if hasW == true {
		delete(r.workerConnections, c.Id)
	}
}

// gateway 链接
func (r *RegisterServerEvent) gatewayConnect(c *network.TcpServerClient, msg RegisterMessage) {
	if msg.Address == "" {
		println("address not found")
		_ = c.Close()
		return
	}
	// 推入列表
	r.gatewayConnections[c.Id] = msg.Address
	r.broadcastAddresses(0)
}

// worker 链接
func (r *RegisterServerEvent) workerConnect(c *network.TcpServerClient, msg RegisterMessage) {
	// 推入列表
	r.workerConnections[c.Id] = c
	r.broadcastAddresses(0)
}

/*
向 BusinessWorker 广播 gateway 内部通讯地址
0 全部发生
 */
func (r *RegisterServerEvent) broadcastAddresses(id uint32) {
	type ConList struct {
		Event     string   `json:"event"`
		Addresses []string `json:"addresses"`
	}
	data := ConList{Event: "broadcast_addresses"}

	for _, address := range r.gatewayConnections {
		data.Addresses = append(data.Addresses, address)
	}

	jsonByte, _ := json.Marshal(data)
	sendMsg := string(jsonByte)

	if id != 0 {
		worker := r.workerConnections[id]
		worker.Send(sendMsg)
		return
	}

	for _, worker := range r.workerConnections {
		worker.Send(sendMsg)
	}
}
