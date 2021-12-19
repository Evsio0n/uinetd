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

var defaultConf = "/etc/uinetd.conf"

//ipv4RegexPattern matches ipv4 string like "127.0.0.1"
var ipv4RegexPattern = `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`

//ipv6onlyRegexPattern matching string like "2001:da8::"
var ipv6onlyRegexPattern = "^[0-9a-fA-F]*:([0-9a-fA-F]|:)*$"

// ipv6WithSquareBracketPattern matching string like "[2001:da8::]"
var ipv6WithSquareBracketPattern = "^\\[[0-9a-fA-F]*:([0-9a-fA-F]|:)*\\]$"

//portRegexPattern matches port string like "8080"
var portRegexPattern = `^([0-9]|[1-9][0-9]|[1-9][0-9]{2}|[1-9][0-9]{3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`

func ReadConf() {
	//read conf file
	file, err := os.OpenFile(defaultConf, os.O_RDONLY, 0666)
	if err != nil {
		logger.NewError(fmt.Sprintf("open conf file error:%e", err))
		return
	}
	defer file.Close()

	//read per line and parse
	buffer := bufio.NewReader(file)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		//parse line
		parseLine(line)
	}
}

//parseLine parse line and set config
func parseLine(str string) {
	//deal with comment
	if strings.HasPrefix(str, "#") {
		//do nothing
	} else {
		//separate line
		lineprefix := strings.Fields(str)
		//check line prefix
		if len(lineprefix) == 0 {
			//do nothing
		} else {
			//check line prefix
			switch lineprefix[0] {
			//TODO: adding deny and allow chain
			case "deny":

			case "allow":

			case "loglevel":
				atoi, err := strconv.Atoi(lineprefix[1])
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
				//check if have 5 fields per line
				if len(lineprefix) == 5 {
					//check if both address and port is valid
					if check.IP(lineprefix[0]) && check.IP(lineprefix[2]) {
						//check port is valid
						if check.Port(lineprefix[1]) && check.Port(lineprefix[3]) {
							//check protocol
							if check.Protocol(lineprefix[4]) {
								//TODO: add to config

							} else {
								logger.NewError(fmt.Sprintf("protocol: \"%s\" is invalid.", lineprefix[4]))
							}
						} else {
							logger.NewError(fmt.Sprintf("port: is invalid.", lineprefix[2]))
						}
					}
				} else {
					logger.NewError(fmt.Sprintf("line: %s is invalid.Skipping", str))
				}

			}
		}
	}
}
