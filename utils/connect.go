package utils

import "net"

type forwardItem struct {
	sourceAddr string
	sourcePort int
	destAddr   string
	destPort   int
}

type forwardItems []forwardItem

type connectController struct {
	//connectionType is the type of connection
	//0:raw connect
	//1:tcp
	//2:udp
	//3:tcp+udp
	connectionType int

	//tcp connection
	tcpConn *net.TCPConn

	//udp connection
	udpConn *net.UDPConn

	//raw connection
	rawConn *net.Conn
}

//forwardItemController is a controller for forwardItem.
type forwardItemController struct {
	item       forwardItem
	controller *connectController
}

type forwardItemControllers []forwardItemController

func (f forwardItemControllers) Len() int {
	return len(f)
}
