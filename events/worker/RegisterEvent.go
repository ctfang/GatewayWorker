package worker

import (
	"encoding/json"
	"github.com/ctfang/GatewayWorker/events"
	"github.com/ctfang/network"
	"github.com/ctfang/network/protocol"
	"github.com/ctfang/network/tcp"
	"log"
)

type RegisterEvent struct {
	listen           network.ListenTcp
	gatewayAddresses map[string]network.Connect
}

type WorkerConnect struct {
	Event     string `json:"event"`
	SecretKey string `json:"secret_key"`
}

type BroadcastAddresses struct {
	Event     string   `json:"event"`
	Addresses []string `json:"addresses"`
}

func (r *RegisterEvent) OnStart(listen network.ListenTcp) {
	log.Println("connect the register to: ", listen.GetAddress().Str)
	r.listen = listen
}

func (*RegisterEvent) OnConnect(c network.Connect) {
	conData := WorkerConnect{
		Event:     "worker_connect",
		SecretKey: events.SecretKey,
	}
	byteStr, _ := json.Marshal(conData)
	go c.Send(string(byteStr))
}

func (r *RegisterEvent) OnMessage(c network.Connect, message interface{}) {
	strMsg := message.(string)
	msgBA := BroadcastAddresses{}
	err := json.Unmarshal([]byte(strMsg), &msgBA)
	if err != nil {
		return
	}
	switch msgBA.Event {
	case "broadcast_addresses":
		for _, addr := range msgBA.Addresses {
			if _, ok := r.gatewayAddresses[addr]; !ok {
				r.gatewayAddresses[addr] = nil
			}
		}
		r.checkGatewayConnections()
	default:
		log.Println("不认识的事件", msgBA)
	}
}

func (*RegisterEvent) OnClose(c network.Connect) {

}

func (*RegisterEvent) OnError(listen network.ListenTcp, err error) {

}

func (r *RegisterEvent) checkGatewayConnections() {
	for addr, con := range r.gatewayAddresses {
		if con == nil {
			worker := tcp.NewClient()
			worker.SetAddress(network.NewAddress(addr))
			worker.SetConnectionEvent(NewGatewayEvent(r, addr))
			worker.SetProtocol(protocol.NewGatewayProtocol())
			go worker.ListenAndServe()
		}
	}
}

/*
连接成功 or 失败
*/
func (r *RegisterEvent) UpdateGatewayConnections(addr string, con network.Connect) {
	if con != nil {
		r.gatewayAddresses[addr] = con
	} else {
		delete(r.gatewayAddresses, addr)
	}
}

/*
连接注册中心
*/
func NewRegisterEvent() network.Event {
	return &RegisterEvent{
		gatewayAddresses: map[string]network.Connect{},
	}
}
