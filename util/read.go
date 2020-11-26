package utils

import (
	"../lib/checkIP"
	"../lib/log"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//Read conf data.

var defaultConf = "/etc/uinetd.conf"

func ReadConfAndDial() error {

	file, err := os.Open(defaultConf)
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
		if len(d) == 0 {
			return fmt.Errorf("%s", "No configuration, exiting...")
		} else {
			if strings.ContainsAny(d[0], "#") {
				// Do nothing with #
			} else if d[1] == "loglevel" {
				intLevel, _ := strconv.Atoi(d[2])
				logger.SetLoglevel(intLevel)
			} else if check.IsIp(d[0]) && check.IsNormalPort(d[1]) && check.IsIp(d[2]) && check.IsNormalPort(d[3]) && check.IsMode(d[5]) {
				//Port must be under 65535
				newForwardItems := forwardingItems{}
				newForwardItems.localIp[i] = d[0]
				newForwardItems.localPorts[i] = atoi(d[1])
				newForwardItems.dstIP[i] = d[2]
				newForwardItems.dstPort[i] = atoi(d[3])
				newForwardItems.method[i] = d[4]
				//Add forwardItems to dial.
				err = dialConn(newForwardItems)
				return err
			}
		}
	}
	file.Close()
	return nil

}
func atoi(str string) int {
	a, _ := strconv.Atoi(str)
	return a
}
