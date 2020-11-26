package check

import (
	"net"
	"strconv"
	"strings"
)

func IsIp(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func IpType(ip string) int {
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	return 0
}

func IsNormalPort(port string) bool {
	//Referring to http://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml, Port must be [0-65535]
	e, err := strconv.Atoi(port)
	if err != nil {
		return false
	} else if e < 65535 && e > 0 {
		return true
	}
	return false
}

func IsMode(str string) bool {
	strLow := strings.ToLower(str)
	if strLow == "udp" || strLow == "tcp" || strLow == "all" {
		return true
	}
	return false
}
