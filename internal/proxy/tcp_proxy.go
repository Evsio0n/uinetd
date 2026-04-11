package proxy

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/evsio0n/uinetd/internal/config"
	"github.com/evsio0n/uinetd/internal/logger"
)

// TCPProxy TCP代理服务器
type TCPProxy struct {
	rule     config.ForwardRule
	cfg      *config.Config
	logger   *logger.Logger
	listener net.Listener
	wg       sync.WaitGroup // 追踪所有活跃连接
}

// NewTCPProxy 创建新的TCP代理
func NewTCPProxy(rule config.ForwardRule, cfg *config.Config, log *logger.Logger) *TCPProxy {
	return &TCPProxy{
		rule:   rule,
		cfg:    cfg,
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
	p.listener = listener

	p.logger.LogInfo("TCP 代理启动: %s -> %s:%d",
		listenAddr, p.rule.ConnectAddress, p.rule.ConnectPort)

	return nil
}

// Run 运行 TCP 代理的 accept 循环，直到 ctx 取消或 listener 关闭
func (p *TCPProxy) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		if p.listener != nil {
			p.listener.Close()
		}
	}()

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			// 检查是否是因为 listener 被关闭（正常退出）
			select {
			case <-ctx.Done():
				return
			default:
			}
			p.logger.LogError("接受连接失败: %v", err)
			continue
		}

		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.handleConnection(conn)
		}()
	}
}

// Shutdown 优雅关闭 TCP 代理，等待所有活跃连接完成
func (p *TCPProxy) Shutdown(timeout time.Duration) {
	if p.listener != nil {
		p.listener.Close()
	}

	// 等待所有活跃连接完成，或超时
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		p.logger.LogInfo("TCP 代理已优雅关闭 (所有连接已完成)")
	case <-time.After(timeout):
		p.logger.LogInfo("TCP 代理关闭超时，强制退出 (等待了 %v)", timeout)
	}
}

// handleConnection 处理TCP连接
func (p *TCPProxy) handleConnection(clientConn net.Conn) {
	defer func() {
		if err := clientConn.Close(); err != nil {
			p.logger.LogDebug("关闭客户端连接失败: %v", err)
		}
	}()

	// 获取客户端信息
	var clientAddr *net.TCPAddr
	if addr, ok := clientConn.RemoteAddr().(*net.TCPAddr); ok {
		clientAddr = addr
	} else {
		p.logger.LogError("无法解析客户端地址类型: %T", clientConn.RemoteAddr())
		return
	}

	// 设置 TCP KeepAlive（仅在配置了 keepalive 时启用）
	if p.cfg.KeepAlive > 0 {
		if tcpConn, ok := clientConn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(p.cfg.KeepAlive)
		}
	}

	// 连接到目标服务器
	targetAddr := net.JoinHostPort(p.rule.ConnectAddress, strconv.Itoa(p.rule.ConnectPort))
	var targetConn net.Conn
	var err error
	if p.cfg.Timeout > 0 {
		targetConn, err = net.DialTimeout("tcp", targetAddr, p.cfg.Timeout)
	} else {
		targetConn, err = net.Dial("tcp", targetAddr)
	}
	if err != nil {
		p.logger.LogError("连接目标服务器失败 %s: %v", targetAddr, err)
		return
	}
	defer func() {
		if err := targetConn.Close(); err != nil {
			p.logger.LogDebug("关闭目标连接失败: %v", err)
		}
	}()

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
