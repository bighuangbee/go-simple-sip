package tools

import (
	"strconv"
	"strings"
)

func ParseAddr(addr string)(ip string, port uint16){
	index := strings.Index(addr, ":")
	if index < 10{
		return "", 0
	}
	p, _ := strconv.Atoi(addr[index+1:])
	return addr[:index], uint16(p)
}
