package logger

import (
	logger "log"
	"os"
)

var defaultLogPath = "/var/log/uinetd.log"

var logLevelConfig int

func SetLoglevel(loglevel int) {
	switch loglevel {
	case 1:
		logLevelConfig = 1
	case 2:
		logLevelConfig = 2
	case 3:
		logLevelConfig = 3
	case 4:
		logLevelConfig = 4
	default:
		logLevelConfig = 0
	}
}

func InitialLog() error {
	f, err := os.OpenFile(defaultLogPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	if _, err := f.Write([]byte("appended some data\n")); err != nil {
		f.Close()
		// ignore error; Write error takes precedence
		logger.Fatal(err)
	}
	if f.Close() != nil {
		logger.Fatal(err)
		return err
	}
	return nil
}

func logNew(logPath string, loglevel int) {

}
