package logger

import "log/syslog"

var defaultLogPath = "/var/log/uinetd.log"

var logLevelConfig int
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
		logLevelConfig = 0
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
	if logLevelConfig <= 3 {
		sysLogger, err = syslog.New(syslog.LOG_SYSLOG, "uinetd")
		if err != nil {
			return err
		} else {
			sysLogger, err = syslog.New(syslog.LOG_SYSLOG, "uinetd")
			if err != nil {
				return err
			}
		}
		sysLogger.Info("")
	}
	return nil
}
func NewInfo(str string) {
	sysLogger.Info(str)
}

func NewError(str string) {
	sysLogger.Err(str)
}
