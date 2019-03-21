package gateway

import (
	"GatewayWorker/network"
	"errors"
)

type WorkerRouter struct {
	// worker ConnectionId 映射 worker
	workers map[uint32]network.Connect
	// client ConnectionId 映射 client
	Clients map[uint32]network.Connect
	// client ConnectionId 映射 worker ，记录已经通讯过的通道 client =》 worker
	clientList map[uint32]network.Connect
}

var Router = &WorkerRouter{
	workers:    map[uint32]network.Connect{},
	clientList: map[uint32]network.Connect{},
	Clients:    map[uint32]network.Connect{},
}

func (w *WorkerRouter) GetWorker(c network.Connect) (network.Connect, error) {
	cid := c.GetConnectionId()
	if worker, ok := w.clientList[cid]; ok {
		// 已经通信过的通道
		return worker, nil
	} else {
		// 随机分配一个worker
		for _, worker := range w.workers {
			w.clientList[cid] = worker
			return worker, nil
		}
	}
	// 不存在
	return nil, errors.New("找不到worker")

}

func (w *WorkerRouter) AddedWorker(worker network.Connect) {
	w.workers[worker.GetConnectionId()] = worker
}

/**
worker 断开
*/
func (w *WorkerRouter) DeleteWorker(worker network.Connect) {
	cid := worker.GetConnectionId()
	delete(w.workers, cid)
	for clientId, worker := range w.clientList {
		if cid == worker.GetConnectionId() {
			delete(w.clientList, clientId)
		}
	}
}

/**
新增客户端，并且建立路由映射
*/
func (w *WorkerRouter) AddedClient(c network.Connect) {
	ConnectionId := c.GetConnectionId()
	if _, ok := w.Clients[ConnectionId]; ok {
		w.DeleteClient(ConnectionId)
	}

	w.Clients[ConnectionId] = c
}

func (w *WorkerRouter) GetClient(ConnectionId uint32) (network.Connect, error) {
	c, ok := w.Clients[ConnectionId]
	if ok {
		return c, nil
	}

	return nil, errors.New("客户端不存在")
}

/**
删除 客户端
*/
func (w *WorkerRouter) DeleteClient(ConnectionId uint32) {
	delete(w.clientList, ConnectionId)
	delete(w.Clients, ConnectionId)
}
