package tcp

import (
	"GatewayWorker/network"
	"GatewayWorker/network/protocol"
	"log"
	"net"
)

type Server struct {
	address  *network.Address
	event    network.Event
	protocol network.Protocol
	lastId   uint32
}

func NewServer() network.ListenTcp {
	return &Server{
		lastId: 0,
	}
}

/*
设置监听地址
*/
func (server *Server) SetAddress(address *network.Address) {
	server.address = address
}

func (server *Server) GetAddress() *network.Address {
	if server.address == nil {
		panic("没有设置地址")
	}
	return server.address
}

/*
设置信息事件
*/
func (server *Server) SetConnectionEvent(event network.Event) {
	server.event = event
}

func (server *Server) GetConnectionEvent() network.Event {
	if server.event == nil {
		panic("没有设置事件")
	}
	return server.event
}

/*
设置协议解析格式
*/
func (server *Server) SetProtocol(protocol network.Protocol) {
	server.protocol = protocol
}

/*
获取协议解析
*/
func (server *Server) GetProtocol() network.Protocol {
	if server.protocol == nil {
		server.SetProtocol(protocol.NewNewline())
	}
	return server.protocol
}

/*
启动监听
*/
func (server *Server) ListenAndServe() {
	address := server.GetAddress()
	event := server.GetConnectionEvent()
	listener, err := net.Listen("tcp", address.Str)
	if err != nil {
		go event.OnError(server, &ListenError{address})
		log.Fatal("Error starting TCP server.", address.Str)
		return
	}
	defer listener.Close()
	go event.OnStart(server)
	for {
		con, _ := listener.Accept()
		server.lastId += 1
		go server.newConnection(con)
	}
}

/*
新的连接
*/
func (server *Server) newConnection(con net.Conn) {
	var connection network.Connect = NewConnection(con, server, server.lastId)
	event := server.GetConnectionEvent()

	defer event.OnClose(connection)
	go event.OnConnect(connection)

	for {
		message, err := connection.Read()
		if err != nil {
			break
		}
		go event.OnMessage(connection, message)
	}
}