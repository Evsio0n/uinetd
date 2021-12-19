package main

import (
	"os"
)

func initial() error {
	return nil
}

//TODO: Communicate with dbus or systemd (Only systemd can do start or stop project.)
func main() {
	err := initial()
	if err != nil {
		os.Exit(1)
	}
}
