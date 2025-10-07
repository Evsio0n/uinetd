package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
}

// ParseConfig 解析配置文件
func ParseConfig(filename string) (*Config, error) {
	// #nosec G304: 配置文件路径由用户/调用方提供，属预期可变输入
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer func() { _ = file.Close() }()

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
				if err == nil && level >= 1 && level <= 4 {
					config.LogLevel = level
				}
			}
			continue
		}

		// 解析转发规则
		parts := strings.Fields(line)
		if len(parts) < 5 {
			fmt.Printf("警告: 第 %d 行格式不正确，已忽略: %s\n", lineNum, line)
			continue
		}

		bindPort, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("警告: 第 %d 行端口格式不正确: %s\n", lineNum, parts[1])
			continue
		}

		connectPort, err := strconv.Atoi(parts[3])
		if err != nil {
			fmt.Printf("警告: 第 %d 行端口格式不正确: %s\n", lineNum, parts[3])
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

// stripBrackets 移除 IPv6 地址周围的方括号
func stripBrackets(addr string) string {
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		return addr[1 : len(addr)-1]
	}
	return addr
}
