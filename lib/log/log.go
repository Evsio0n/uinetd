package log

import (
	logger "log"
	"os"
)

var defaultLogPath = "/var/log/uinetd.log"

func logNew(logPath string) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	if _, err := f.Write([]byte("appended some data\n")); err != nil {
		f.Close()

		// ignore error; Write error takes precedence
		logger.Fatal(err)
	}
	if f.Close() != nil {
		logger.Fatal(err)
	}
}
