package main

import (
	"./lib/log"
	"bufio"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

var nullpd string
var bindAddress string
var bindPort string
var bindToAddress string
var bindToPort string
var t int = 0
var q int = 0

type c2 net.Listener

var checkBool int = 0

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

func loadConfigurationFile() {
	filename := "/etc/uinetd.conf"
	f, err := os.Open(filename)
	if err != nil {
		log.Error("Error at open the configuration file.")
		os.Exit(1)
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
	s1 := strings.Replace(readString, "  ", " ", -1)
	regStr := "\\s{2,}"
	reg, _ := regexp.Compile(regStr)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spcIndex := reg.FindStringIndex(string(s2))
	for len(spcIndex) > 0 {
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(s2))
	}
	i := len(string(s2))
	for c := 0; c < i; c++ {
		if string(s2[c]) == " " {
			t++
		}
		switch t {
		case 1:
			if checkBool == 0 {
				bindAddress = string(s2[q:c])
				q = c + 1
				checkBool = 1
			}

		case 2:
			if checkBool == 1 {
				bindPort = string(s2[q:c])
				q = c + 1
				checkBool = 2
			}
		case 3:
			if checkBool == 2 {
				bindToAddress = string(s2[q:c])
				bindToPort = string(s2[c+1 : i])
				checkBool = 3
				q = c
			}
		case 4:
			log.Error("Set pram wrong!")
			os.Exit(1)
		}
	}
	if checkIPAddressType(bindToAddress) && checkIPAddressType(bindAddress) {
		log.Info("Verified!")
	}
}

func checkIPAddressType(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}
