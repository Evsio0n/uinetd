package utils

import (
	"bufio"
	"fmt"
	"github.com/evsio0n/uinetd/pkg/check"
	"github.com/evsio0n/uinetd/pkg/logger"
	"os"
	"strconv"
	"strings"
)

//Read conf data.

var defaultConf = "/etc/uinetd/uinetd.conf"

func ReadConf() error {
	//read conf file
	file, err := os.OpenFile(defaultConf, os.O_RDONLY, 0666)
	if err != nil {
		logger.NewError(fmt.Sprintf("open conf file error:%s", err.Error()))
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.NewError(fmt.Sprintf("Error:%s", err.Error()))
		}
	}(file)

	//read per line and parse
	phraseLine := 1
	buffer := bufio.NewReader(file)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		//parse line
		parseLine(line, phraseLine)
		phraseLine++
	}
	//When finish reading we should let all ForwardItems Up.
	//Goto Startup
	ForwardItems.Start()
	return nil
}

//parseLine parse line and set config, then add items to ForwardItems.
func parseLine(str string, phraseLine int) {
	//deal with comment
	if strings.HasPrefix(str, "#") {
		//do nothing
	} else {
		//separate line
		linePrefix := strings.Fields(str)
		//check line prefix
		if len(linePrefix) == 0 {
			//do nothing
		} else {
			//check line prefix
			switch linePrefix[0] {
			//TODO: adding deny and allow chain
			case "deny":

			case "allow":

			case "loglevel":
				atoi, err := strconv.Atoi(linePrefix[1])
				if err != nil {
					logger.NewError(fmt.Sprintf("parse loglevel error:%e", err))
				} else {
					if atoi > 0 && atoi < 5 {
						logger.SetLoglevel(atoi)
					} else {
						logger.NewError(fmt.Sprintf("loglevel: %d is invalid.", atoi))
					}
				}
			default:
				//check if we have 5 fields per line
				if len(linePrefix) == 5 {
					//check if both address and port is valid
					if check.IP(linePrefix[0]) && check.IP(linePrefix[2]) {
						//check port is valid
						if check.Port(linePrefix[1]) && check.Port(linePrefix[3]) {
							//check protocol
							if check.Protocol(linePrefix[4]) {
								//TODO: add to config
								ForwardItems.Add(linePrefix[0], safeAtoI(linePrefix[1]), linePrefix[2], safeAtoI(linePrefix[3]), linePrefix[4])
							} else {
								logger.NewError(fmt.Sprintf("protocol: \"%s\" is invalid in line:%d,Skipping...", linePrefix[4], phraseLine))
							}
						} else {
							logger.NewError(fmt.Sprintf("port is invalid in line:%d,Skipping...", phraseLine))
						}
					}
				} else {
					logger.NewError(fmt.Sprintf("line:%d is invalid.Skipping...", phraseLine))
				}

			}
		}
	}
}

//safeAtoI is format string to int. Use it when safe.
func safeAtoI(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
