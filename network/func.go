package network

import (
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"strconv"
)

// 构建分布式唯一id
func Bin2hex(ip string, port uint16, id uint32) string {
	var msgByte []byte
	ipUint32 := Ip2long(ip)
	var buf32 = make([]byte, 4)
	var bug16 = make([]byte, 2)
	binary.BigEndian.PutUint32(buf32, ipUint32)
	msgByte = append(msgByte, buf32...)
	binary.BigEndian.PutUint16(bug16, port)
	msgByte = append(msgByte, bug16...)
	binary.BigEndian.PutUint32(buf32, id)
	msgByte = append(msgByte, buf32...)
	return hex.EncodeToString(msgByte)
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

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return ip
}
