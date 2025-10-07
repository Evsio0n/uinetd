package logger

import (
	"fmt"
	"log"
	"time"
)

// Logger 日志记录器
type Logger struct {
	Level int
}

// NewLogger 创建新的日志记录器
func NewLogger(level int) *Logger {
	return &Logger{Level: level}
}

// LogError 记录错误 (所有级别都会记录)
func (l *Logger) LogError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// LogProhibited 记录被禁止的连接 (级别 2, 3, 4)
func (l *Logger) LogProhibited(fromIP string, fromPort int, toIP string, toPort int) {
	if l.Level >= 2 {
		if l.Level >= 3 {
			log.Printf("[PROHIBITED] %s - %s:%d -> %s:%d",
				time.Now().Format("2006-01-02 15:04:05"),
				fromIP, fromPort, toIP, toPort)
		} else {
			log.Printf("[PROHIBITED] %s", time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}

// LogConnection 记录连接 (级别 4)
func (l *Logger) LogConnection(fromIP string, fromPort int, toIP string, toPort int) {
	if l.Level >= 4 {
		log.Printf("[CONNECTION] %s - %s:%d -> %s:%d",
			time.Now().Format("2006-01-02 15:04:05"),
			fromIP, fromPort, toIP, toPort)
	}
}

// LogInfo 记录信息
func (l *Logger) LogInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// LogDebug 记录调试信息
func (l *Logger) LogDebug(format string, v ...interface{}) {
	if l.Level >= 4 {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// GetTimePrefix 获取时间前缀
func GetTimePrefix() string {
	return fmt.Sprintf("[%s]", time.Now().Format("2006-01-02 15:04:05"))
}
