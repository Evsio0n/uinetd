package logger

import (
	"log/syslog"
)

var defaultLogPath = "/var/log/uinetd.log"

var logLevelConfig = 1
var sysLogger *syslog.Writer
var err error

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
		logLevelConfig = 1
	}
}

//This job should be done by syslog
/*func InitialLog() error {
	os.Create(defaultLogPath)
	f, err := os.OpenFile(defaultLogPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Fatal(err)
		return err
	}
	//create file if not exist
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
}*/

//InitialLog initial log with syslog
func InitialLog() error {
	//default set to lazy output before from SetLogLevel.
	logLevelConfig = 0
	sysLogger, err = syslog.New(syslog.LOG_SYSLOG, "uinetd")
	if err != nil {
		return err
	}
	sysLogger.Info("Starting uinetd from now.")
	return nil
}
func NewInfo(str string) {
	if checkLogLevelIsVerbose() {
		sysLogger.Info(str)
	}
}

func NewError(str string) {
	sysLogger.Err(str)
}

func checkLogLevelIsVerbose() bool {
	//if logLevelConfig is 1 then only output level ERROR to syslog
	if logLevelConfig != 1 {
		return true
	}
	return false
}
