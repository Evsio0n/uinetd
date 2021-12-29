package utils

import (
	"fmt"
	"github.com/evsio0n/uinetd/pkg/logger"
	"io"
	"net"
)

type forwardItem struct {
	sourceAddr  string
	sourcePort  int
	destAddr    string
	destPort    int
	forwardType string
}

type forwardItems []forwardItem
type connectController struct {
	//connectionType is the type of connection
	//0:raw connect
	//1:tcp
	//2:udp
	//3:tcp+udp
	connectionType int

	//tcp listener
	tcpListener *net.TCPListener
	//tcp listener connection
	tcpListenerConn *net.TCPConn
	//tcp  destination connection
	tcpDstConn *net.TCPConn

	//udp connection
	udpConn *net.UDPConn

	//udp listener
	udpListener *net.UDPConn

	//raw connection
	rawConn *net.Conn
	//raw listener
	rawListener *net.Listener
}

//forwardItemController is a controller for forwardItem.
type forwardItemController struct {
	item       forwardItem
	controller *connectController
}

type forwardItemControllers []forwardItemController

var ForwardItems forwardItems
var ForwardItemControllers forwardItemControllers

//Add a forwardItem to ForwardItems.
func (f *forwardItems) Add(sourceAddr string, sourcePort int, destAddr string, destPort int, forwardType string) {
	*f = append(*f, forwardItem{sourceAddr, sourcePort, destAddr, destPort, forwardType})
	ForwardItemControllers.Add(forwardItem{
		sourceAddr:  sourceAddr,
		sourcePort:  sourcePort,
		destAddr:    destAddr,
		destPort:    destPort,
		forwardType: forwardType,
	})

}

func (f *forwardItems) Start() {

}

func (f *forwardItems) Del(num int) {

}

//Stop all forwardItem.
func (f *forwardItemController) Stop() {

}

//Flush  all forwardItems.
func (f *forwardItems) Flush() {
	//Delete all forwardItems.
	*f = nil
}

func (c *forwardItemControllers) Add(item forwardItem) {
	*c = append(*c, forwardItemController{item, &connectController{}})
}

//DialTcp is a function to dial tcp with forwardItemController.
func (c *forwardItemController) DialTcp() {
	//Binding local address.
	localAddr, err := net.ResolveTCPAddr("tcp", c.item.sourceAddr+":"+string(c.item.sourcePort))
	if err != nil {
		logger.NewError(fmt.Sprintf("ResolveTCPAddr error: %s", err.Error()))
	} else {
		//Listen on local address.
		c.controller.tcpListener, err = net.ListenTCP("tcp", localAddr)
		if err != nil {
			logger.NewError(fmt.Sprintf("ListenTCP error: %s", err.Error()))
		} else {
			remoteAddr, err := net.ResolveTCPAddr("tcp", c.item.destAddr+":"+string(c.item.destPort))
			if err != nil {
				logger.NewError(fmt.Sprintf("ResolveTCPAddr error: %s", err.Error()))
			} else {
				go func() {
					for {
						//Accept a connection.
						c.controller.tcpListenerConn, err = c.controller.tcpListener.AcceptTCP()
						if err != nil {
							logger.NewError(fmt.Sprintf("AcceptTCP error: %s", err.Error()))
						}
						//Dial to remote address.
						c.controller.tcpDstConn, err = net.DialTCP("tcp", nil, remoteAddr)
						if err != nil {
							logger.NewError(fmt.Sprintf("DialTCP error: %s", err.Error()))
						}
						//setting timeout for both connection to reduce io wait.

						//Copy data from local connection to remote connection.
						go copyData(c.controller.tcpListenerConn, c.controller.tcpDstConn)
						go copyData(c.controller.tcpDstConn, c.controller.tcpListenerConn)

					}
				}()
			}
		}
	}
}

func copyData(src, dst net.Conn) {
	defer func(src net.Conn) {
		err := src.Close()
		if err != nil {
			logger.NewInfo(fmt.Sprintf("Close error: %s", err.Error()))
		}
	}(src)
	defer func(dst net.Conn) {
		err := dst.Close()
		if err != nil {
			logger.NewInfo(fmt.Sprintf("Close error: %s", err.Error()))
		}
	}(dst)
	_, err := io.Copy(dst, src)
	if err != nil {
		logger.NewInfo(fmt.Sprintf("Copy error: %s", err.Error()))
	}

}

//forwardItemControllerChan is a channel for forwardItemController.
type forwardItemControllerChan chan forwardItemController

var ForwardItemControllerChan forwardItemControllerChan
