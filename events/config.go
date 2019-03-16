package events

import (
	"GatewayWorker/network"
	"GatewayWorker/network/tcp"
)

var GatewayAddress *network.Address
var WorkerAddress *network.Address
var RegisterAddress *network.Address
var SecretKey string

var test = tcp.NewServer()
