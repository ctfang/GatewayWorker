package events

type LogicEvent interface {
	OnStart()
	// 新链接
	OnConnect(clientId string)
	// 新信息
	OnMessage(clientId string, body []byte)
	// 链接关闭
	OnClose(clientId string)
}
