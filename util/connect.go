package utila

import (
	"fmt"
	"io"
	"net"
)

type forwardingItems struct {
	localIp    []string
	localPorts []int
	method     []string
	dstIP      []string
	dstPort    []int
}
type listenType struct {
	TcpError     error
	UdpError     error
	tcpListener  *net.TCPListener
	udpConn      *net.UDPConn
	port         int
	listenerType string
	dialIP       string
	dialPort     int
}

type listenGroup chan listenType

func dialConn(item forwardingItems) error {
	for i := 0; i < len(item.method); i++ {
		switch item.method[i] {
		case "tcp":
			dialTcpConn(item.localIp[i], item.localPorts[i], item.dstIP[i], item.dstPort[i])
			return nil
		case "udp":
			dialUdpConn(item.localIp[i], item.localPorts[i], item.dstIP[i], item.dstPort[i])
			return nil
		case "all":
			dialAllConn(item.localIp[i], item.localPorts[i], item.dstIP[i], item.dstPort[i])
			return nil
		}
	}
	return fmt.Errorf("%s", "Error")
}

func dialTcpConn(ip string, ports int, dstIP string, dstPort int) {
	l := listenType{}
	addr := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: ports,
	}
	listenerTcp, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		l.TcpError = err
		l.tcpListener = nil
	} else {
		client, err := listenerTcp.Accept()
		if err != nil {

		}
		target, _ := net.Dial("tcp", dstIP+":"+string(ports))
		defer target.Close()
		defer client.Close()
		go func() { io.Copy(target, client) }()
		go func() { io.Copy(client, target) }()
		l.tcpListener = listenerTcp
		l.listenerType = "tcp"
		l.port = ports
		l.dialIP = dstIP
		l.dialPort = dstPort
		l.addToChannel()
	}
	//TODO: ignoring something error and add to log level 1

}

func dialUdpConn(ip string, ports int, dstIP string, dstPort int) {
	l := listenType{}
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: ports,
	}
	listenerUdp, err := net.ListenUDP("tcp", &addr)
	if err != nil {
		l.tcpListener = nil
	} else {
		l.UdpError = err
		l.udpConn = listenerUdp
		l.listenerType = "tcp"
		l.port = ports
		l.dialIP = dstIP
		l.dialPort = dstPort
		l.addToChannel()
	}
	//TODO: ignoring something error and add to log level 1
}

func dialAllConn(ip string, ports int, dstIP string, dstPort int) {
	l := listenType{}
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: ports,
	}
	addrs := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: ports,
	}
	listenerUdp, err := net.ListenUDP("udp", &addr)
	listenerTcp, errs := net.ListenTCP("tcp", &addrs)
	if err != nil && errs != nil {
		l.tcpListener = nil
		l.UdpError = err
		l.TcpError = err
	} else {
		l.udpConn = listenerUdp
		l.tcpListener = listenerTcp
		l.listenerType = "all"
		l.port = ports
		l.dialIP = dstIP
		l.dialPort = dstPort
		l.addToChannel()
	}
	//TODO: ignoring something error and add to log level 1
}

func (l listenType) addToChannel() {

}

func tcpTransfer() {

}
