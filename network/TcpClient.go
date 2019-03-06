package network

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strconv"
)

type TcpEventInterface interface {
	OnConnect(clint *TcpServerConnection)
	OnMessage(clint *TcpServerConnection, message interface{})
	OnClose(clint *TcpServerConnection)
}

type TcpServerConnection struct {
	Addr     *Address
	conn     net.Conn
	Protocol Protocol
	Event    TcpEventInterface
	SendChan chan []byte
}

func (t *TcpServerConnection) Connect()error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.Addr.Str);
	conn, err := net.DialTCP("tcp", nil,tcpAddr)

	if err != nil {
		log.Printf("connect failed, err : %v\n", err.Error())
		return err
	}
	t.conn = conn

	go t.Event.OnConnect(t)
	go t.readPump()
	go t.writePump()
	return nil
}

func (t *TcpServerConnection) writePump() {
	defer t.Close()
	for {
		select {
		case text := <-t.SendChan:
			_, err := t.conn.Write(text)
			if err != nil {
				log.Println("Error writing to stream.",err)
				break
			}
		}
	}
}

func (t *TcpServerConnection) readPump() {
	reader := bufio.NewReader(t.conn)
	for {
		message, err := t.Protocol.ReadString(reader)
		if err != nil {
			t.Close()
			break
		}
		go t.Event.OnMessage(t, message)
	}
}

func (t *TcpServerConnection) Send(str string) {
	t.SendChan <- []byte(t.Protocol.WriteString(str))
}

func (t *TcpServerConnection) Close() {
	t.conn.Close()
	t.Event.OnClose(t)
}

func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1>255 || ip2>255 || ip3>255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return ip
}
