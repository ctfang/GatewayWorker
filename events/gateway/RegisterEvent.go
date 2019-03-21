package gateway

import (
	"GatewayWorker/events"
	"GatewayWorker/network"
	"encoding/json"
	"log"
	"time"
)

type ConList struct {
	Event     string `json:"event"`
	Address   string `json:"address"`
	SecretKey string `json:"secret_key"`
}

type RegisterEvent struct {
	retry  int16
	listen network.ListenTcp
}

// @error
func (r *RegisterEvent) OnError(listen network.ListenTcp, err error) {
	r.retry++
	log.Println("注册中心连接失败，2秒后重试", r.retry)
	ticker := time.NewTicker(time.Second * 2)
	select {
	case <-ticker.C:
		listen.ListenAndServe()
		break
	}
}

func (r *RegisterEvent) OnStart(listen network.ListenTcp) {
	log.Println("connect the register to: ", listen.GetAddress().Str)
	r.listen = listen
}

func (*RegisterEvent) OnConnect(c network.Connect) {
	conData := ConList{
		Event:     "gateway_connect",
		Address:   events.WorkerAddress.Str,
		SecretKey: events.SecretKey,
	}
	byteStr, _ := json.Marshal(conData)
	go c.Send(string(byteStr))
}

func (*RegisterEvent) OnMessage(c network.Connect, message interface{}) {
	log.Println("gateway 收到注册中心的信息 ", message)
}

func (r *RegisterEvent) OnClose(c network.Connect) {
	// 注册中心 关闭,定时检查
	log.Print("注册中心断开连接，2秒后重连 ", r.retry)
	ticker := time.NewTicker(time.Second * 2)
	select {
	case <-ticker.C:
		r.listen.ListenAndServe()
	}
}

// 连接注册中心
func NewRegisterEvent() network.Event {
	return &RegisterEvent{}
}
