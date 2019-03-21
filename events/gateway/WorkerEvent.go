package gateway

import (
	"GatewayWorker/network"
	"GatewayWorker/network/protocol"
	"log"
)

/*
接受worker发上来的数据
*/
type WorkerEvent struct {
	HandleFunc map[uint8]func(c network.Connect, message protocol.GatewayMessage)
}

func (*WorkerEvent) OnError(listen network.ListenTcp, err error) {

}

func (w *WorkerEvent) OnStart(listen network.ListenTcp) {
	log.Println("worker server listening at: ", listen.GetAddress().Str)

	if w.HandleFunc == nil {
		w.HandleFunc = map[uint8]func(c network.Connect, message protocol.GatewayMessage){}
		WorkerHandle := NewWorkerHandle()

		w.HandleFunc[protocol.CMD_WORKER_CONNECT] = WorkerHandle.OnWorkerConnect
		w.HandleFunc[protocol.CMD_GATEWAY_CLIENT_CONNECT] = WorkerHandle.OnGatewayClientConnect
		w.HandleFunc[protocol.CMD_SEND_TO_ONE] = WorkerHandle.OnSendToOne
		w.HandleFunc[protocol.CMD_KICK] = WorkerHandle.OnSendToOne
		w.HandleFunc[protocol.CMD_DESTROY] = WorkerHandle.OnDestroy
		w.HandleFunc[protocol.CMD_SEND_TO_ALL] = WorkerHandle.OnSendToAll
		w.HandleFunc[protocol.CMD_SELECT] = WorkerHandle.OnSelect

		w.HandleFunc[protocol.CMD_GET_GROUP_ID_LIST] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_SET_SESSION] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_UPDATE_SESSION] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_GET_SESSION_BY_CLIENT_ID] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_GET_ALL_CLIENT_SESSIONS] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_IS_ONLINE] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_BIND_UID] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_SEND_TO_UID] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_LEAVE_GROUP] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_UNGROUP] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_SEND_TO_GROUP] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_GET_CLIENT_SESSIONS_BY_GROUP] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_GET_CLIENT_COUNT_BY_GROUP] = WorkerHandle.OnTodo
		w.HandleFunc[protocol.CMD_GET_CLIENT_ID_BY_UID] = WorkerHandle.OnTodo
	}
}

func (*WorkerEvent) OnConnect(c network.Connect) {

}

func (w *WorkerEvent) OnMessage(c network.Connect, message interface{}) {
	msg := message.(protocol.GatewayMessage)

	if handle, ok := w.HandleFunc[msg.Cmd]; ok {
		handle(c, msg)
	} else {
		log.Println("不认识的命令", msg.Cmd, message)
	}
}

func (*WorkerEvent) OnClose(c network.Connect) {
	c.Close()
	Router.DeleteWorker(c)
}

func NewWorkerEvent() network.Event {
	return &WorkerEvent{}
}
