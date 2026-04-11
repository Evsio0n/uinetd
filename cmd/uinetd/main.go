package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/evsio0n/uinetd/internal/config"
	"github.com/evsio0n/uinetd/internal/logger"
	"github.com/evsio0n/uinetd/internal/proxy"
)

const (
	defaultConfigFile  = "/etc/uinetd.conf"
	shutdownTimeout    = 15 * time.Second
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

// ProxyManager 管理所有代理实例的生命周期
type ProxyManager struct {
	tcpProxies []*proxy.TCPProxy
	udpProxies []*proxy.UDPProxy
}

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

	// 创建代理管理器
	pm := &ProxyManager{}

	// 启动所有代理
	for _, rule := range cfg.Rules {
		if err := pm.startProxy(rule, cfg, log); err != nil {
			log.LogError("启动代理失败: %v", err)
			// 启动失败时，关闭已经启动的代理
			pm.shutdownAll(log)
			os.Exit(1)
		}
	}

	log.LogInfo("所有代理服务已启动")

	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 运行所有代理的 accept/read 循环
	for _, p := range pm.tcpProxies {
		go p.Run(ctx)
	}
	for _, p := range pm.udpProxies {
		go p.Run(ctx)
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.LogInfo("收到信号 %v, 正在优雅关闭...", sig)

	// 取消 context，通知所有 goroutine 停止
	cancel()

	// 优雅关闭所有代理
	pm.shutdownAll(log)

	log.LogInfo("uinetd 已停止")
}

// startProxy 根据协议类型启动相应的代理
func (pm *ProxyManager) startProxy(rule config.ForwardRule, cfg *config.Config, log *logger.Logger) error {
	protocol := strings.ToLower(rule.Protocol)

	switch protocol {
	case "tcp":
		p := proxy.NewTCPProxy(rule, cfg, log)
		if err := p.Start(); err != nil {
			return err
		}
		pm.tcpProxies = append(pm.tcpProxies, p)
		return nil

	case "udp":
		p := proxy.NewUDPProxy(rule, cfg, log)
		if err := p.Start(); err != nil {
			return err
		}
		pm.udpProxies = append(pm.udpProxies, p)
		return nil

	case "all", "both":
		tcpP := proxy.NewTCPProxy(rule, cfg, log)
		if err := tcpP.Start(); err != nil {
			return err
		}
		pm.tcpProxies = append(pm.tcpProxies, tcpP)

		udpP := proxy.NewUDPProxy(rule, cfg, log)
		if err := udpP.Start(); err != nil {
			// UDP 启动失败，关闭已启动的 TCP
			tcpP.Shutdown(shutdownTimeout)
			return err
		}
		pm.udpProxies = append(pm.udpProxies, udpP)

		return nil

	case "raw":
		log.LogInfo("RAW 协议暂不支持，规则已忽略: %s:%d", rule.BindAddress, rule.BindPort)
		return nil

	default:
		return fmt.Errorf("不支持的协议类型: %s", rule.Protocol)
	}
}

// shutdownAll 优雅关闭所有代理
func (pm *ProxyManager) shutdownAll(log *logger.Logger) {
	for _, p := range pm.tcpProxies {
		p.Shutdown(shutdownTimeout)
	}
	for _, p := range pm.udpProxies {
		p.Shutdown(shutdownTimeout)
	}
}
