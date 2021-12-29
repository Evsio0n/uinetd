package main

import (
	"github.com/evsio0n/uinetd/pkg/logger"
	"github.com/evsio0n/uinetd/utils"
	"os"
)

func initial() error {
	err := logger.InitialLog()
	if err != nil {
		return err
	}
	err = utils.ReadConf()
	if err != nil {
		return err
	}
	return nil
}

//TODO: Communicate with dbus or systemd (Only systemd can do start or stop project.)
func main() {
	err := initial()
	if err != nil {
		os.Exit(1)
	}
}
