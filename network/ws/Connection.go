package ws

import (
	"GatewayWorker/network"
	"bytes"
	"github.com/gorilla/websocket"
	"net"
	"regexp"
	"strconv"
	"time"
)

type Connection struct {
	cid    uint32
	uid    string
	con    *websocket.Conn
	pro    network.Protocol
	extend interface{}
}

func NewConnection(con *websocket.Conn, server network.ListenTcp, cid uint32) network.Connect {
	return &Connection{
		cid: cid,
		con: con,
		pro: server.GetProtocol(),
	}
}

func (c *Connection) GetConnectionId() uint32 {
	return c.cid
}

func (c *Connection) SetUid(uid string) {
	c.uid = uid
}

func (c *Connection) GetUid() string {
	return c.uid
}

func (c *Connection) Send(msg interface{}) bool {
	_ = c.con.SetWriteDeadline(time.Now().Add(writeWait))
	message := c.pro.Write(c, msg)
	_ = c.con.WriteMessage(websocket.TextMessage, message)
	return true
}

func (c *Connection) GetCon() net.Conn {
	panic("websocket不能获取原连接")
}

func (c *Connection) Close() {
	c.con.Close()
}

func (c *Connection) Read() (interface{}, error) {
	_, message, err := c.con.ReadMessage()
	if err != nil {
		return message, err
	}
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	return message, nil
}

func (c *Connection) SetExtend(extend interface{}) {
	c.extend = extend
}

func (c *Connection) GetExtend() interface{} {
	return c.extend
}

func (c *Connection) GetIp() string {
	ipStr := c.con.RemoteAddr().String()
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

func (c *Connection) GetPort() uint16 {
	ipStr := c.con.RemoteAddr().String()
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
	port := uint16(temp)
	return port
}
