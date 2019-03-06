package network

import (
	"github.com/gorilla/websocket"
	"net"
	"regexp"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// websocket client
type TcpClientConnection struct {
	// 进程内唯一，递增
	id     uint32
	// 连接唯一，分布
	uid    string
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	// 扩展对象
	Extend   interface{}
}

func (c *TcpClientConnection) Close() {
	c.conn.Close()
}

func (c *TcpClientConnection) Send(meg []byte) {
	c.send <- meg
}

// 进程内id
func (c *TcpClientConnection) GetConnectionId() uint32 {
	return c.id
}

// 唯一id
func (c *TcpClientConnection) SetUid(str string) {
	c.uid = str
}

// 唯一id
func (c *TcpClientConnection) GetUid() string {
	return c.uid
}

func (c *TcpClientConnection) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TcpClientConnection) GetPort() (port uint16) {
	ipStr := c.conn.RemoteAddr().String()
	r := `\:(\d{1,5})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return 0
	}
	ips := reg.FindStringSubmatch(ipStr)
	if ips == nil {
		return 0
	}
	temp, _ := strconv.Atoi(ips[1])
	port = uint16(temp)
	return
}

func (c *TcpClientConnection) GetIp() (ip string) {
	ipStr := c.conn.RemoteAddr().String()
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return ""
	}
	ips := reg.FindStringSubmatch(ipStr)
	if ips == nil {
		return ""
	}

	return ips[0]
}