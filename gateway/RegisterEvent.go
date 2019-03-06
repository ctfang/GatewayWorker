package gateway

import (
	"GoGatewayWorker/network"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var GatewayAddress *network.Address
var WorkerAddress *network.Address
var RegisterAddress *network.Address
var SecretKey string

/*
连接注册中心客户端
 */
type RegisterEvent struct {
	// 内部通讯地址
}

type ConList struct {
	Event     string `json:"event"`
	Address   string `json:"address"`
	SecretKey string `json:"secret_key"`
}

func (r *RegisterEvent) OnConnect(clint *network.TcpServerConnection) {
	conData := ConList{
		Event:     "gateway_connect",
		Address:   WorkerAddress.Str,
		SecretKey: SecretKey,
	}
	byteStr, _ := json.Marshal(conData)
	go clint.Send(string(byteStr))
	log.Println("已经连接注册中心", clint.Addr)
}

func (*RegisterEvent) OnMessage(clint *network.TcpServerConnection, message interface{}) {
	fmt.Println(message)
}

// 关闭
func (*RegisterEvent) OnClose(clint *network.TcpServerConnection) {
	// 注册中心 关闭,定时检查
	log.Print("注册中心 关闭,定时检查")
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			err := clint.Connect()
			if err == nil {
				ticker.Stop()
				return
			}
		}
	}
}
