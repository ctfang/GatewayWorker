package network

import (
	"strconv"
	"strings"
)

type Address struct {
	Str  string
	Ip   string
	Port uint16
}

func NewAddress(addr string) *Address {
	strS := strings.Split(addr, ":")
	if len(strS) != 2 {
		panic("格式错误")
	}
	Port, _ := strconv.ParseInt(strS[1], 10, 64)
	return &Address{
		Str:  addr,
		Ip:   strS[0],
		Port: uint16(Port),
	}
}
