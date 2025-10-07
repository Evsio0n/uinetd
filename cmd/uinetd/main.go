package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"uinetd/internal/config"
	"uinetd/internal/logger"
	"uinetd/internal/proxy"
)

const (
	defaultConfigFile = "/etc/uinetd.conf"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

func main() {
	// 命令行参数
	configFile := flag.String("c", defaultConfigFile, "配置文件路径")
	showVersion := flag.Bool("v", false, "显示版本信息")
	flag.Parse()

	if *showVersion {
		fmt.Printf("uinetd version %s\n", version)
		fmt.Printf("构建时间: %s\n", buildTime)
		return
	}

	// 解析配置文件
	cfg, err := config.ParseConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Rules) == 0 {
		fmt.Fprintf(os.Stderr, "警告: 没有找到任何转发规则\n")
		os.Exit(1)
	}

	// 创建日志记录器
	log := logger.NewLogger(cfg.LogLevel)
	log.LogInfo("uinetd v%s 启动中...", version)
	log.LogInfo("日志级别: %d", cfg.LogLevel)
	log.LogInfo("加载了 %d 条转发规则", len(cfg.Rules))

	// 启动所有代理
	for _, rule := range cfg.Rules {
		if err := startProxy(rule, log); err != nil {
			log.LogError("启动代理失败: %v", err)
			os.Exit(1)
		}
	}

	log.LogInfo("所有代理服务已启动")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.LogInfo("收到信号 %v, 正在关闭...", sig)
	log.LogInfo("uinetd 已停止")
}

// startProxy 根据协议类型启动相应的代理
func startProxy(rule config.ForwardRule, log *logger.Logger) error {
	protocol := strings.ToLower(rule.Protocol)

	switch protocol {
	case "tcp":
		p := proxy.NewTCPProxy(rule, log)
		return p.Start()

	case "udp":
		p := proxy.NewUDPProxy(rule, log)
		return p.Start()

	case "all", "both":
		// 同时启动 TCP 和 UDP
		tcpProxy := proxy.NewTCPProxy(rule, log)
		if err := tcpProxy.Start(); err != nil {
			return err
		}

		udpProxy := proxy.NewUDPProxy(rule, log)
		if err := udpProxy.Start(); err != nil {
			return err
		}

		return nil

	case "raw":
		log.LogInfo("RAW 协议暂不支持，规则已忽略: %s:%d", rule.BindAddress, rule.BindPort)
		return nil

	default:
		return fmt.Errorf("不支持的协议类型: %s", rule.Protocol)
	}
}
