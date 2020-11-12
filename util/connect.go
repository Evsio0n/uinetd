package util

import (
	"fmt"
	"net"
)

type forwardingItems struct {
	ip     []string
	ports  [][]int
	method []string
}
type listenType struct {
	tcpListener  []*net.TCPListener
	udpConn      []*net.UDPConn
	ports        []int
	listenerType string
}

type listenGroup chan listenType

func dialConn(item forwardingItems) (err error) {
	for i := 0; i < len(item.ip); i++ {
		switch item.method[i] {
		case "tcp":
			callTcpConn(item.ip[i], item.ports[i])
			return nil
		case "udp":
			callUdpConn(item.ip[i], item.ports[i])
			return nil
		case "all":
			callAllConn(item.ip[i], item.ports[i])
			return nil
		default:
			return fmt.Errorf("")
		}
	}
	return fmt.Errorf("")
}

func callTcpConn(ip string, ports []int) {
	l := listenType{}
	for i := 0; i < len(ports); i++ {
		addr := net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: ports[i],
		}
		listenerTcp, err := net.ListenTCP("tcp", &addr)
		if err == nil {
			l.tcpListener[i] = listenerTcp
			l.listenerType = "tcp"
			l.ports = ports
			l.addToChannel()
		}
		//TODO: ignoring something error and add to log level 1
	}
}

func callUdpConn(ip string, ports []int) {
	l := listenType{}
	for i := 0; i < len(ports); i++ {
		addrUdp := net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: ports[i],
		}
		listenerUdp, err := net.ListenUDP("udp", &addrUdp)
		if err == nil {
			l.udpConn[i] = listenerUdp
			l.listenerType = "udp"
			l.ports = ports
			l.addToChannel()
		}
	}
}

func callAllConn(ip string, ports []int) {
	l := listenType{}
	for i := 0; i < len(ports); i++ {
		addrUdp := net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: ports[i],
		}
		addrTcp := net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: ports[i],
		}
		listener, err := net.ListenUDP("udp", &addrUdp)
		if err == nil {
			l.udpConn[i] = listener
		}
		tcpListener, err := net.ListenTCP("tcp", &addrTcp)
		if err == nil {
			l.tcpListener[i] = tcpListener
			l.listenerType = "all"
			l.ports = ports
			l.addToChannel()
		}
	}
}

func (l listenType) addToChannel() {

}
