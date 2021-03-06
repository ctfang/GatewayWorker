package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 发给worker，gateway有一个新的连接
const CMD_ON_CONNECT = 1;

// 发给worker的，客户端有消息
const CMD_ON_MESSAGE = 3;

// 发给worker上的关闭链接事件
const CMD_ON_CLOSE = 4;

// 发给gateway的向单个用户发送数据
const CMD_SEND_TO_ONE = 5;

// 发给gateway的向所有用户发送数据
const CMD_SEND_TO_ALL = 6;

// 发给gateway的踢出用户
// 1、如果有待发消息，将在发送完后立即销毁用户连接
// 2、如果无待发消息，将立即销毁用户连接
const CMD_KICK = 7;

// 发给gateway的立即销毁用户连接
const CMD_DESTROY = 8;

// 发给gateway，通知用户session更新
const CMD_UPDATE_SESSION = 9;

// 获取在线状态
const CMD_GET_ALL_CLIENT_SESSIONS = 10;

// 判断是否在线
const CMD_IS_ONLINE = 11;

// client_id绑定到uid
const CMD_BIND_UID = 12;

// 解绑
const CMD_UNBIND_UID = 13;

// 向uid发送数据
const CMD_SEND_TO_UID = 14;

// 根据uid获取绑定的clientid
const CMD_GET_CLIENT_ID_BY_UID = 15;

// 加入组
const CMD_JOIN_GROUP = 20;

// 离开组
const CMD_LEAVE_GROUP = 21;

// 向组成员发消息
const CMD_SEND_TO_GROUP = 22;

// 获取组成员
const CMD_GET_CLIENT_SESSIONS_BY_GROUP = 23;

// 获取组在线连接数
const CMD_GET_CLIENT_COUNT_BY_GROUP = 24;

// 按照条件查找
const CMD_SELECT = 25;

// 获取在线的群组ID
const CMD_GET_GROUP_ID_LIST = 26;

// 取消分组
const CMD_UNGROUP = 27;

// worker连接gateway事件
const CMD_WORKER_CONNECT = 200;

// 心跳
const CMD_PING = 201;

// GatewayClient连接gateway事件
const CMD_GATEWAY_CLIENT_CONNECT = 202;

// 根据client_id获取session
const CMD_GET_SESSION_BY_CLIENT_ID = 203;

// 发给gateway，覆盖session
const CMD_SET_SESSION = 204;

// 当websocket握手时触发，只有websocket协议支持此命令字
const CMD_ON_WEBSOCKET_CONNECT = 205;

// 包体是标量
const FLAG_BODY_IS_SCALAR = 0x01;

// 通知gateway在send时不调用协议encode方法，在广播组播时提升性能
const FLAG_NOT_CALL_ENCODE = 0x02;


type GatewayProtocol struct {
	cacheLen    int
	cacheReader []byte
}
// "Npack_len/Ccmd/Nlocal_ip/nlocal_port/Nclient_ip/nclient_port/Nconnection_id/Cflag/ngateway_port/Next_len"
type GatewayMessage struct {
	PackageLen   uint32
	Cmd          uint8
	LocalIp      uint32
	LocalPort    uint16
	ClientIp     uint32
	ClientPort   uint16
	ConnectionId uint32
	Flag         uint8
	GatewayPort  uint16

	ExtLen  uint32
	ExtData string
	Body    string
}

func (t *GatewayProtocol) ReadString(reader *bufio.Reader) (interface{}, error) {
	if t.cacheLen > 0 {
		data := t.cacheReader
		dataLen := t.cacheLen
		t.cacheLen = 0
		t.cacheReader = []byte{}

		// 包长度,不足重新读取
		if dataLen < 4 {
			dataTem := make([]byte, 1024)
			dataLenTem, err := reader.Read(dataTem)
			if err != nil {
				return nil, err
			}
			dataLen = dataLen + dataLenTem
			data = append(dataTem)
		}
		// 包长度大于实际数据长度，还有数据未传完整
		PackageLen := uint32(binary.BigEndian.Uint32(data[0:4]))
		intPackageLen := int(PackageLen)
		if intPackageLen > dataLen {
			// 实际数据不足
			t.cacheLen = dataLen
			t.cacheReader = data
			return t.ReadString(reader)
		} else if intPackageLen == dataLen {
			return t.ReadStruct(data), nil
		} else if intPackageLen < dataLen {
			// 实际比需要的长
			t.cacheLen = dataLen - intPackageLen
			t.cacheReader = data[intPackageLen:]
			return t.ReadStruct(data), nil
		}
	} else {
		data := make([]byte, 1024)

		dataLen, err := reader.Read(data)
		if err != nil {
			return nil, err
		}
		// 包长度,不足直接抛弃连接
		if dataLen < 4 {
			return nil, err
		}
		// 包长度大于实际数据长度，还有数据未传完整
		PackageLen := uint32(binary.BigEndian.Uint32(data[0:4]))
		intPackageLen := int(PackageLen)
		if intPackageLen > dataLen {
			// 实际数据不足
			t.cacheLen = dataLen
			t.cacheReader = data
			return t.ReadString(reader)
		} else if intPackageLen == dataLen {
			return t.ReadStruct(data), nil
		} else if intPackageLen < dataLen {
			// 实际比需要的长
			t.cacheLen = dataLen - intPackageLen
			t.cacheReader = data[intPackageLen:]
			return t.ReadStruct(data), nil
		}
	}
	return nil, errors.New("不能解析")
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func (t *GatewayProtocol) WriteString(msg interface{}) []byte {
	var msgByte []byte

	value := reflect.ValueOf(msg)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.String:
			var bufStr []byte
			bufStr = []byte(field.String())
			msgByte = append(msgByte, bufStr...)
		case reflect.Uint8:
			msgByte = append(msgByte, uint8(field.Uint()))
		case reflect.Uint16:
			var buf16 = make([]byte, 2)
			binary.BigEndian.PutUint16(buf16, uint16(field.Uint()))
			msgByte = append(msgByte, buf16...)
		case reflect.Uint32:
			var buf32 = make([]byte, 4)
			binary.BigEndian.PutUint32(buf32, uint32(field.Uint()))
			msgByte = append(msgByte, buf32...)
		default:
			fmt.Println("不知道的类型",field.Type())
		}
	}

	return msgByte
}

func (t *GatewayProtocol) ReadStruct(data []byte) GatewayMessage {
	Message := GatewayMessage{
		PackageLen:   uint32(binary.BigEndian.Uint32(data[0:4])),
		Cmd:          data[4],
		LocalIp:      uint32(binary.BigEndian.Uint32(data[5:9])),
		LocalPort:    uint16(binary.BigEndian.Uint16(data[9:11])),
		ClientIp:     uint32(binary.BigEndian.Uint32(data[11:15])),
		ClientPort:   uint16(binary.BigEndian.Uint16(data[15:17])),
		ConnectionId: uint32(binary.BigEndian.Uint32(data[17:21])),
		Flag:         data[21],
		GatewayPort:  uint16(binary.BigEndian.Uint16(data[22:24])),
		ExtLen:       uint32(binary.BigEndian.Uint32(data[24:28])),
	}
	Message.ExtData = string(data[28 : 28+Message.ExtLen])
	Message.Body = string(data[(28+Message.ExtLen) : (Message.PackageLen)])

	return Message
}
