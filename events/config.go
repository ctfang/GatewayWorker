package events

import (
	"GatewayWorker/network"
)

var GatewayAddress *network.Address
var WorkerAddress *network.Address
var RegisterAddress *network.Address
var SecretKey string

var BussinessEvent LogicEvent
