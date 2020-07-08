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
var bindAddress string
var bindPort string
var bindToAddress string
var bindToPort string
var intTest bool = false
var t int = 0
var q int = 0

type c2 net.Listener

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

func createTcpThread(readAddress string) c2 {
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
	defer f.Close()
	r := bufio.NewReader(f)
	var check bool = false
	var buf []byte
	for {
		buf, err = r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
		} else if strings.Contains("\n", string(buf)) {
			check = true
		} else {
			check = false
		}
		if check {
			//Do nothing
		} else {
			if strings.Contains(string(buf), "#") {
			} else {
				readString := string(buf)
				readIt(readString)
			}
		}
	}
}

func readIt(readString string) {
	i := len(readString)
	for c := 0; c == i-1; c++ {
		if string(readString[c]) == " " {
			t++
			q++
			//TODO: Do some thing to read the port.
			if string(readString[c]) == string(readString[c-1]) {
				t--
				//Some people didn't see the fucking instruction, so we should delete deprecated space
			}
		}
		switch t {
		case 1:
			//Read bind address
			bindAddress = string(readString[q:c])
			log.Info(bindAddress)
			q = c
		case 2:
			//Read bind port
			bindPort = string(readString[q:c])
			log.Info(bindPort)
			q = c
		case 3:
			//Read bind to address
			bindToAddress = string(readString[q:c])
			log.Info(bindToAddress)
			q = c
		case 4:
			//Read bind to port
			bindToPort = string(readString[q:c])
			log.Info(bindToPort)
			q = c

		}

	}
}

func checkIPAddressType(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}
