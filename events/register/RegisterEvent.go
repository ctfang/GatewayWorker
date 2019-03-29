package register

import (
	"GatewayWorker/network"
	"encoding/json"
	"fmt"
	"log"
)

type RegisterMessage struct {
	Event     string `json:"event"`
	Address   string `json:"address"`
	SecretKey string `json:"secret_key"`
}

type RegisterEvent struct {
	// ID 地址保存
	gatewayConnections map[uint32]string
	workerConnections  map[uint32]network.Connect
	// 秘要
	SecretKey string
}

func (*RegisterEvent) OnStart(listen network.ListenTcp) {
	log.Println("register server listening at: ", listen.GetAddress().Str)
}

func (*RegisterEvent) OnConnect(c network.Connect) {

}

func (r *RegisterEvent) OnMessage(c network.Connect, message interface{}) {
	var data RegisterMessage
	err := json.Unmarshal([]byte(message.(string)), &data)
	if err != nil {
		fmt.Println(err)
		c.Close()
		return
	}
	if r.SecretKey != "" {
		if data.SecretKey != r.SecretKey {
			fmt.Println("秘要不对")
			c.Close()
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
		c.Close()
	}
}

func (r *RegisterEvent) OnClose(c network.Connect) {
	_, hasG := r.gatewayConnections[c.GetConnectionId()]
	if hasG == true {
		delete(r.gatewayConnections, c.GetConnectionId())
		r.broadcastAddresses(0)
	}

	_, hasW := r.workerConnections[c.GetConnectionId()]
	if hasW == true {
		delete(r.workerConnections, c.GetConnectionId())
	}
}

func (*RegisterEvent) OnError(listen network.ListenTcp, err error) {
	log.Println("注册中心启动失败", err)
}

// gateway 链接
func (r *RegisterEvent) gatewayConnect(c network.Connect, msg RegisterMessage) {
	if msg.Address == "" {
		println("address not found")
		c.Close()
		return
	}
	// 推入列表
	r.gatewayConnections[c.GetConnectionId()] = msg.Address
	r.broadcastAddresses(0)
}

// worker 链接
func (r *RegisterEvent) workerConnect(c network.Connect, msg RegisterMessage) {
	// 推入列表
	r.workerConnections[c.GetConnectionId()] = c
	r.broadcastAddresses(0)
}

/*
向 BusinessWorker 广播 gateway 内部通讯地址
0 全部发生
*/
func (r *RegisterEvent) broadcastAddresses(id uint32) {
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

func NewRegisterEvent() network.Event {
	return &RegisterEvent{
		gatewayConnections: make(map[uint32]string),
		workerConnections:  make(map[uint32]network.Connect),
	}
}
