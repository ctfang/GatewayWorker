package network

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WebSocket struct {
	Addr     *Address
	upgrader websocket.Upgrader
	Event    WsEventInterface
	// 最新值id
	clientId uint32
}

type WsEventInterface interface {
	// 启动前
	OnStart()
	// 有客户端连接
	OnConnect(clint *TcpClientConnection)
	// 有信息
	OnMessage(clint *TcpClientConnection, message []byte)
	// 关闭
	OnClose(clint *TcpClientConnection)
}

// 开启监听
func (ws *WebSocket) ListenAndServe() {
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	go ws.Event.OnStart()

	http.HandleFunc("/", ws.Upgrade)
	err := http.ListenAndServe(ws.Addr.Str, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 生产新id
func (w *WebSocket) MakeId() uint32 {
	if w.clientId >= MaxClient {
		w.clientId = 0
	}
	w.clientId++
	return w.clientId
}

func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &TcpClientConnection{
		conn: conn,
		send: make(chan []byte, 256),
		id:   ws.MakeId(),
	}

	go ws.Event.OnConnect(client)
	go ws.writePump(client)
	go ws.readPump(client)
}

func (ws *WebSocket) writePump(c *TcpClientConnection) {
	defer ws.Event.OnClose(c)

	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (ws *WebSocket) readPump(c *TcpClientConnection) {
	defer ws.Event.OnClose(c)

	// 信息size上限
	c.conn.SetReadLimit(maxMessageSize)
	// 设置底层网络连接的读取截止日期。读取超时后，websocket连接状态已损坏，所有将来的读取都将返回错误。t的零值意味着读取不会超时。
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// Pong 信息
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			//	log.Printf("error: %v", err)
			//}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		go ws.Event.OnMessage(c, message)
	}
}
