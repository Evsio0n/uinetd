package utila

import (
	check "../lib/checkIP"
	logger "../lib/log"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//TODO:Read conf data.

var defaultConf = "/etc/uintd.conf"

func readConf(strConf string) error {
	var conf string
	if strConf == "" {
		conf = defaultConf
	} else {
		conf = strConf
	}
	file, err := os.Open(conf)
	if err != nil {
		file.Close()
		return fmt.Errorf("%s:%e", "Error when loading file ", err)
	}
	br := bufio.NewReader(file)
	for {
		i := 0
		i++
		a, _, c := br.ReadLine()
		//Read configure per line.
		if c == io.EOF {
			break
		}
		d := strings.Fields(string(a))
		if strings.ContainsAny(d[1], "#") {
			// Do nothing with #
		} else if d[1] == "loglevel" {
			intLevel, _ := strconv.Atoi(d[2])
			logger.SetLoglevel(intLevel)
		} else if check.IsIp(d[1]) && check.IsNormalPort(d[2]) && check.IsIp(d[3]) && check.IsNormalPort(d[4]) && check.CheckMode(d[5]) {
			//Port must be under 65535
			newForwardItems := forwardingItems{}
			newForwardItems.localIp[i] = d[1]
			newForwardItems.localPorts[i] = atoi(d[2])
			newForwardItems.dstIP[i] = d[3]
			newForwardItems.dstPort[i] = atoi(d[4])
			newForwardItems.method[i] = d[5]
			dialConn(newForwardItems)
		}
	}
	file.Close()
	return nil

}
func atoi(str string) int {
	a, _ := strconv.Atoi(str)
	return a
}
