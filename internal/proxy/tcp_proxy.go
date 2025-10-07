package proxy

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"uinetd/internal/config"
	"uinetd/internal/logger"
)

// TCPProxy TCP代理服务器
type TCPProxy struct {
	rule   config.ForwardRule
	logger *logger.Logger
}

// NewTCPProxy 创建新的TCP代理
func NewTCPProxy(rule config.ForwardRule, log *logger.Logger) *TCPProxy {
	return &TCPProxy{
		rule:   rule,
		logger: log,
	}
}

// Start 启动TCP代理服务器
func (p *TCPProxy) Start() error {
	listenAddr := net.JoinHostPort(p.rule.BindAddress, strconv.Itoa(p.rule.BindPort))

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("无法监听 TCP %s: %v", listenAddr, err)
	}

	p.logger.LogInfo("TCP 代理启动: %s -> %s:%d",
		listenAddr, p.rule.ConnectAddress, p.rule.ConnectPort)

	go func() {
		defer func() { _ = listener.Close() }()
		for {
			conn, err := listener.Accept()
			if err != nil {
				p.logger.LogError("接受连接失败: %v", err)
				continue
			}

			go p.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection 处理TCP连接
func (p *TCPProxy) handleConnection(clientConn net.Conn) {
	defer func() { _ = clientConn.Close() }()

	// 获取客户端信息
	var clientAddr *net.TCPAddr
	if addr, ok := clientConn.RemoteAddr().(*net.TCPAddr); ok {
		clientAddr = addr
	} else {
		p.logger.LogError("无法解析客户端地址类型: %T", clientConn.RemoteAddr())
		return
	}

	// 连接到目标服务器
	targetAddr := net.JoinHostPort(p.rule.ConnectAddress, strconv.Itoa(p.rule.ConnectPort))
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		p.logger.LogError("连接目标服务器失败 %s: %v", targetAddr, err)
		return
	}
	defer func() { _ = targetConn.Close() }()

	// 记录连接
	p.logger.LogConnection(
		clientAddr.IP.String(),
		clientAddr.Port,
		p.rule.ConnectAddress,
		p.rule.ConnectPort,
	)

	// 双向转发数据
	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 -> 目标服务器
	go func() {
		defer wg.Done()
		if _, err := io.Copy(targetConn, clientConn); err != nil {
			p.logger.LogDebug("转发客户端->目标失败: %v", err)
		}
		if tcp, ok := targetConn.(*net.TCPConn); ok {
			if err := tcp.CloseWrite(); err != nil {
				p.logger.LogDebug("目标连接半关闭失败: %v", err)
			}
		}
	}()

	// 目标服务器 -> 客户端
	go func() {
		defer wg.Done()
		if _, err := io.Copy(clientConn, targetConn); err != nil {
			p.logger.LogDebug("转发目标->客户端失败: %v", err)
		}
		if tcp, ok := clientConn.(*net.TCPConn); ok {
			if err := tcp.CloseWrite(); err != nil {
				p.logger.LogDebug("客户端连接半关闭失败: %v", err)
			}
		}
	}()

	wg.Wait()
}
