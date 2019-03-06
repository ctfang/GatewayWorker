package network

import "bufio"

type Protocol interface{
	// 读入处理
	ReadString(reader *bufio.Reader) (interface{}, error)
	// 发送处理
	WriteString(msg interface{}) []byte
}
