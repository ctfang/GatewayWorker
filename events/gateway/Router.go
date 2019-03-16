package gateway

import (
	"GatewayWorker/network"
	"errors"
)

type WorkerRouter struct {
	// 消费者 worker 注册列表
	workers map[uint32]network.Connect

	// client 列表
	Clients map[uint32]network.Connect
	// client ConnectionId 映射 worker id
	clientList map[uint32]uint32
}

var Router = &WorkerRouter{
	workers:    map[uint32]network.Connect{},
	clientList: map[uint32]uint32{},
	Clients:    map[uint32]network.Connect{},
}

func (w *WorkerRouter) GetWorker(c network.Connect) network.Connect {
	workerId := w.clientList[c.GetConnectionId()]
	return w.workers[workerId]
}

func (w *WorkerRouter) AddedWorker(worker network.Connect) {
	w.workers[worker.GetConnectionId()] = worker
}

/**
worker 断开
*/
func (w *WorkerRouter) DeleteWorker(worker network.Connect) {
	delete(w.workers, worker.GetConnectionId())
	for clientId, workerId := range w.clientList {
		if workerId == worker.GetConnectionId() {
			delete(w.clientList, clientId)
		}
	}
}

/**
新增客户端，并且建立路由映射
*/
func (w *WorkerRouter) AddedClient(c network.Connect) (uint32, error) {
	for workerId, _ := range w.workers {
		w.clientList[c.GetConnectionId()] = workerId
		w.Clients[c.GetConnectionId()] = c
		return workerId, nil
	}
	var err = errors.New("未有worker连接")
	return uint32(0), err
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
