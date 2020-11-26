package main

import (
	"./lib/logger"
	"./util"
	"fmt"
	logs "github.com/evsio0n/log"
	"os"
)

func initial() error {
	err := logger.InitialLog()
	if err != nil {
		return fmt.Errorf("%s:%e", "Error when opening log file at /var/log/uinetd.log, Error message:", err)
	} else {
		err = utils.ReadConfAndDial()
		if err != nil {
			return fmt.Errorf("%s:%e", "Error when read configuration file, Error message:", err)
		}
	}
	return nil
}

//TODO: Communicate with dbus or systemd (Only systemd can do start or stop project.)
func main() {
	err := initial()
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
}
