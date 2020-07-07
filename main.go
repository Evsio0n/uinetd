package main

import (
	"./lib/log"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var nullpd string

func readyToCreateThread(readAddress string, readProtocol string) {
	tcpd := "tcp"
	udpd := "udp"
	if readProtocol == tcpd {
		nullpd = tcpd
		createTcpThread(readAddress)
	} else if readProtocol == udpd {
		nullpd = udpd
		createUcpThread(readAddress)
	} else if readProtocol == "all" || readProtocol == "All" || readProtocol == "ALL" {
		//For some people who didn't see the fucking instruction.
		nullpd = "all"
		createAllThread(readAddress)
	} else {
		panic("We shouldn't be there.")
	}
}

type c net.Listener

func createTcpThread(readAddress string) c {
	c, ers := net.Listen("tcp", readAddress)
	if ers != nil {
		log.Error(ers)
	}
	return c
}

func createUcpThread(readAddress string) {
	_, ers := net.Listen("tcp", readAddress)
	if ers != nil {
		log.Error(ers)
	}
}

func createAllThread(readAddress string) {
	createUcpThread(readAddress)
	createTcpThread(readAddress)
}

func main() {
	filename := "/etc/uinetd.conf"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	//关闭文件
	defer f.Close()
	r := bufio.NewReader(f)
	//新建一个缓冲区，把内容先放在缓冲区
	var check bool = false
	var buf []byte
	for {
		//遇到'\n'结束读取, 但是'\n'也读取进入
		buf, err = r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF { //文件已经结束
				break
			}
			log.Error(err)
		} else if strings.Contains(string(buf), "\n") {
			check = true
		} else {
			check = false
		}
		if check {
			//Do nothing
		} else {
			if strings.Contains(string(buf), "#") {
			} else {
				fmt.Print(string(buf))
			}
		}
	}
}
