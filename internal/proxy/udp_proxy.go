package proxy

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"uinetd/internal/config"
	"uinetd/internal/logger"
)

// UDPProxy UDP代理服务器
type UDPProxy struct {
	rule   config.ForwardRule
	logger *logger.Logger

	// 客户端会话管理
	sessions map[string]*UDPSession
	mu       sync.RWMutex
}

// UDPSession UDP会话
type UDPSession struct {
	clientAddr *net.UDPAddr
	targetConn *net.UDPConn
	lastActive time.Time
}

// NewUDPProxy 创建新的UDP代理
func NewUDPProxy(rule config.ForwardRule, log *logger.Logger) *UDPProxy {
	return &UDPProxy{
		rule:     rule,
		logger:   log,
		sessions: make(map[string]*UDPSession),
	}
}

// Start 启动UDP代理服务器
func (p *UDPProxy) Start() error {
	listenAddr := net.JoinHostPort(p.rule.BindAddress, strconv.Itoa(p.rule.BindPort))

	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return fmt.Errorf("解析 UDP 地址失败 %s: %v", listenAddr, err)
	}

	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("无法监听 UDP %s: %v", listenAddr, err)
	}

	p.logger.LogInfo("UDP 代理启动: %s -> %s:%d",
		listenAddr, p.rule.ConnectAddress, p.rule.ConnectPort)

	// 启动会话清理器
	go p.cleanupSessions()

	// 处理接收的数据
	go func() {
		defer func() {
			if err := listener.Close(); err != nil {
				p.logger.LogDebug("关闭 UDP 监听器失败: %v", err)
			}
		}()
		buffer := make([]byte, 65536)

		for {
			n, clientAddr, err := listener.ReadFromUDP(buffer)
			if err != nil {
				p.logger.LogError("UDP 读取失败: %v", err)
				continue
			}

			// 为每个数据包拷贝独立切片，避免共享底层数组被后续读覆盖
			dataCopy := make([]byte, n)
			copy(dataCopy, buffer[:n])
			go p.handlePacket(listener, clientAddr, dataCopy)
		}
	}()

	return nil
}

// handlePacket 处理UDP数据包
func (p *UDPProxy) handlePacket(listener *net.UDPConn, clientAddr *net.UDPAddr, data []byte) {
	sessionKey := clientAddr.String()

	p.mu.Lock()
	session, exists := p.sessions[sessionKey]

	if !exists {
		// 创建新会话
		targetAddr := net.JoinHostPort(p.rule.ConnectAddress, strconv.Itoa(p.rule.ConnectPort))
		udpTargetAddr, err := net.ResolveUDPAddr("udp", targetAddr)
		if err != nil {
			p.logger.LogError("解析目标 UDP 地址失败 %s: %v", targetAddr, err)
			p.mu.Unlock()
			return
		}

		targetConn, err := net.DialUDP("udp", nil, udpTargetAddr)
		if err != nil {
			p.logger.LogError("连接目标 UDP 服务器失败 %s: %v", targetAddr, err)
			p.mu.Unlock()
			return
		}

		session = &UDPSession{
			clientAddr: clientAddr,
			targetConn: targetConn,
			lastActive: time.Now(),
		}
		p.sessions[sessionKey] = session

		// 记录连接
		p.logger.LogConnection(
			clientAddr.IP.String(),
			clientAddr.Port,
			p.rule.ConnectAddress,
			p.rule.ConnectPort,
		)

		// 启动响应接收
		go p.receiveFromTarget(listener, session, sessionKey)
	} else {
		session.lastActive = time.Now()
	}
	p.mu.Unlock()

	// 转发数据到目标服务器
	if _, err := session.targetConn.Write(data); err != nil {
		p.logger.LogDebug("UDP 转发失败: %v", err)
	}
}

// receiveFromTarget 从目标服务器接收响应
func (p *UDPProxy) receiveFromTarget(listener *net.UDPConn, session *UDPSession, sessionKey string) {
	buffer := make([]byte, 65536)

	for {
		if err := session.targetConn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			p.logger.LogDebug("设置 UDP 读超时失败: %v", err)
		}
		n, err := session.targetConn.Read(buffer)
		if err != nil {
			// 超时或错误，关闭会话
			p.mu.Lock()
			delete(p.sessions, sessionKey)
			p.mu.Unlock()
			if err := session.targetConn.Close(); err != nil {
				p.logger.LogDebug("关闭 UDP 目标连接失败: %v", err)
			}
			return
		}

		// 更新活动时间
		p.mu.Lock()
		session.lastActive = time.Now()
		p.mu.Unlock()

		// 转发响应到客户端
		if _, err = listener.WriteToUDP(buffer[:n], session.clientAddr); err != nil {
			p.logger.LogDebug("UDP 响应转发失败: %v", err)
		}
	}
}

// cleanupSessions 清理过期会话
func (p *UDPProxy) cleanupSessions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		p.mu.Lock()
		for key, session := range p.sessions {
			// 5分钟无活动则清理
			if now.Sub(session.lastActive) > 5*time.Minute {
				if err := session.targetConn.Close(); err != nil {
					p.logger.LogDebug("关闭 UDP 目标连接失败: %v", err)
				}
				delete(p.sessions, key)
				p.logger.LogDebug("清理过期 UDP 会话: %s", key)
			}
		}
		p.mu.Unlock()
	}
}
