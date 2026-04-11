package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ForwardRule 表示一条转发规则
type ForwardRule struct {
	BindAddress    string
	BindPort       int
	ConnectAddress string
	ConnectPort    int
	Protocol       string // tcp, udp, all, both, raw
}

// Config 配置结构
type Config struct {
	Rules    []ForwardRule
	LogLevel int

	// KeepAlive TCP KeepAlive 间隔，0 表示使用系统默认值
	KeepAlive time.Duration
	// Timeout 连接目标服务器的超时时间，0 表示不限制
	Timeout time.Duration
}

// ParseConfig 解析配置文件
func ParseConfig(filename string) (*Config, error) {
	// #nosec G304: 配置文件路径由用户/调用方提供，属预期可变输入
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "关闭配置文件失败: %v\n", err)
		}
	}()

	config := &Config{
		Rules:    make([]ForwardRule, 0),
		LogLevel: 1, // 默认日志级别
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 loglevel
		if strings.HasPrefix(line, "loglevel") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				level, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "警告: 第 %d 行日志级别格式不正确: %s\n", lineNum, parts[1])
				} else if level < 1 || level > 4 {
					fmt.Fprintf(os.Stderr, "警告: 第 %d 行日志级别超出范围 (1-4): %d，使用默认值 1\n", lineNum, level)
				} else {
					config.LogLevel = level
				}
			}
			continue
		}

		// 解析 keepalive（秒）
		if strings.HasPrefix(line, "keepalive") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				secs, err := strconv.Atoi(parts[1])
				if err != nil || secs < 0 {
					fmt.Fprintf(os.Stderr, "警告: 第 %d 行 keepalive 格式不正确（需为非负整数）: %s\n", lineNum, parts[1])
				} else {
					config.KeepAlive = time.Duration(secs) * time.Second
				}
			}
			continue
		}

		// 解析 timeout（秒）
		if strings.HasPrefix(line, "timeout") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				secs, err := strconv.Atoi(parts[1])
				if err != nil || secs < 0 {
					fmt.Fprintf(os.Stderr, "警告: 第 %d 行 timeout 格式不正确（需为非负整数）: %s\n", lineNum, parts[1])
				} else {
					config.Timeout = time.Duration(secs) * time.Second
				}
			}
			continue
		}

		// 解析转发规则
		parts := strings.Fields(line)
		if len(parts) < 5 {
			fmt.Fprintf(os.Stderr, "警告: 第 %d 行格式不正确，已忽略: %s\n", lineNum, line)
			continue
		}

		bindPort, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 第 %d 行绑定端口格式不正确: %s\n", lineNum, parts[1])
			continue
		}

		connectPort, err := strconv.Atoi(parts[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 第 %d 行目标端口格式不正确: %s\n", lineNum, parts[3])
			continue
		}

		// 端口范围校验
		if err := validatePort(bindPort, "绑定端口", lineNum); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if err := validatePort(connectPort, "目标端口", lineNum); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		protocol := strings.ToLower(parts[4])

		rule := ForwardRule{
			BindAddress:    stripBrackets(parts[0]),
			BindPort:       bindPort,
			ConnectAddress: stripBrackets(parts[2]),
			ConnectPort:    connectPort,
			Protocol:       protocol,
		}

		config.Rules = append(config.Rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取配置文件时出错: %v", err)
	}

	return config, nil
}

// validatePort 校验端口范围
func validatePort(port int, name string, lineNum int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("警告: 第 %d 行 %s超出范围 (1-65535): %d", lineNum, name, port)
	}
	return nil
}

// stripBrackets 移除 IPv6 地址周围的方括号
func stripBrackets(addr string) string {
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		return addr[1 : len(addr)-1]
	}
	return addr
}
