package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// UINT_MAX
const MaxClient = ^uint32(0)

type TcpServerClient struct {
	Id     uint32
	conn   net.Conn
	Server *TcpServer
	send   chan []byte
}

type TcpServerEvent interface {
	OnStart(tcp *TcpServer)
	// 新链接
	OnConnect(c *TcpServerClient)
	// 新信息
	OnMessage(c *TcpServerClient, message interface{})
	// 链接关闭
	OnClose(c *TcpServerClient)
}

type TcpServer struct {
	// 监听地址
	Address *Address
	// 解析格式
	protocol Protocol
	// 业务处理类
	event TcpServerEvent
	// 最新值id
	clientId uint32
}

// 设置监听地址
func (t *TcpServer) SetAddress(address *Address) {
	t.Address = address
}

// 设置协议解析方式
func (t *TcpServer) SetProtocol(protocol Protocol) {
	t.protocol = protocol
}

func (t *TcpServer) SetEvent(e TcpServerEvent) {
	t.event = e
}

// 开始监听
func (t *TcpServer) ListenAndServe() {
	listener, err := net.Listen("tcp", t.Address.Str)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()
	go t.event.OnStart(t)
	for {
		conn, _ := listener.Accept()

		client := &TcpServerClient{
			Id:     t.MakeId(),
			conn:   conn,
			Server: t,
			send:   make(chan []byte, 256),
		}
		go t.listenClient(client)
	}
}

// 生产新id
func (t *TcpServer) MakeId() uint32 {
	if t.clientId >= MaxClient {
		t.clientId = 0
	}
	t.clientId++
	return t.clientId
}

func (t *TcpServer) listenClient(client *TcpServerClient) {
	// 出发链接事件
	t.event.OnConnect(client)
	go t.writePump(client)
	reader := bufio.NewReader(client.conn)
	for {
		message, err := t.protocol.ReadString(reader)
		if err != nil {
			_ = client.conn.Close()
			t.event.OnClose(client)
			return
		}
		t.event.OnMessage(client, message)
	}
}

func (t *TcpServer) writePump(client *TcpServerClient) {
	for {
		select {
		case message := <-client.send:
			_, err := client.conn.Write(message)
			if err != nil {
				fmt.Println("Error writing to stream.", err)
				break
			}
		}
	}
}

// Send text message to client
func (c *TcpServerClient) Send(message interface{}) {
	byteMes := c.Server.protocol.WriteString(message)
	c.send <- byteMes
}

func (c *TcpServerClient) Close() error {
	return c.conn.Close()
}
