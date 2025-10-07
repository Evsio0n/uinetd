# uinetd - 网络重定向服务器

uinetd 是一个类似于 rinetd 的网络重定向工具，用 Go 语言编写。它可以将 TCP 和 UDP 连接从本地地址和端口重定向到另一个地址和端口。

## 功能特性

- ✅ 支持 TCP 连接转发
- ✅ 支持 UDP 数据包转发
- ✅ 支持同时转发 TCP 和 UDP (ALL/BOTH 协议)
- ✅ 支持 IPv4 和 IPv6
- ✅ 可配置的日志级别
- ✅ 自动 UDP 会话管理和清理
- ✅ 高性能并发处理

## 安装

### 从源代码编译

```bash
# 克隆或下载源代码
cd uinetd

# 编译
go build -o uinetd

# 安装到系统 (可选)
sudo cp uinetd /usr/local/bin/
sudo chmod +x /usr/local/bin/uinetd
```

## 配置

配置文件默认位置：`/etc/uinetd.conf`

### 配置文件格式

```
# 绑定地址    绑定端口    目标地址         目标端口    协议
0.0.0.0      8080       192.168.1.2     80         tcp
0.0.0.0      53         8.8.8.8         53         udp
0.0.0.0      3000       example.com     80         all
[::1]        8080       [2001:DB8::1]   8080       both
```

### 支持的协议

- `tcp` - 仅 TCP 转发
- `udp` - 仅 UDP 转发  
- `all` - 同时转发 TCP 和 UDP
- `both` - 同时转发 TCP 和 UDP (与 all 相同)
- `raw` - 原始套接字 (暂不支持)

### 日志级别

在配置文件中设置：

```
loglevel 4
```

- **Level 1**: 仅记录错误
- **Level 2**: 记录错误和被禁止连接的时间
- **Level 3**: 记录错误和被禁止连接的详细信息
- **Level 4**: 记录所有连接的详细信息和错误

## 使用方法

### 启动服务

```bash
# 使用默认配置文件 /etc/uinetd.conf
sudo uinetd

# 使用自定义配置文件
sudo uinetd -c /path/to/uinetd.conf

# 显示版本信息
uinetd -v
```

### 示例配置

#### 1. HTTP 端口转发

```
# 将本地 8080 端口的 HTTP 流量转发到内网服务器
0.0.0.0      8080       192.168.1.100   80         tcp
```

#### 2. DNS 转发

```
# 将 DNS 查询转发到 Google DNS
0.0.0.0      53         8.8.8.8         53         udp
```

#### 3. 同时转发 TCP 和 UDP

```
# 某些服务同时使用 TCP 和 UDP
0.0.0.0      3000       server.com      3000       all
```

#### 4. IPv6 支持

```
# IPv6 地址转发
[::1]        8080       [2001:DB8::1]   80         tcp
```

## 系统服务配置

### systemd 服务文件

创建 `/etc/systemd/system/uinetd.service`:

```ini
[Unit]
Description=uinetd - Network Redirection Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/uinetd -c /etc/uinetd.conf
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable uinetd
sudo systemctl start uinetd
sudo systemctl status uinetd
```

## 工作原理

### TCP 转发

1. 监听本地指定的 TCP 端口
2. 接受客户端连接
3. 建立到目标服务器的连接
4. 双向转发数据流

### UDP 转发

1. 监听本地指定的 UDP 端口
2. 为每个客户端创建会话
3. 转发数据包到目标服务器
4. 将响应返回给对应的客户端
5. 自动清理过期会话（5分钟无活动）

## 性能优化

- 使用 goroutine 实现高并发
- TCP 连接使用双向流式转发
- UDP 使用会话管理减少连接开销
- 自动清理过期的 UDP 会话

## 安全注意事项

1. **权限**: 监听 1024 以下的端口需要 root 权限
2. **防火墙**: 确保防火墙允许相关端口
3. **访问控制**: 建议配合防火墙规则限制访问
4. **日志**: 在生产环境建议使用较低的日志级别以提高性能

## 故障排查

### 查看日志

```bash
# 如果使用 systemd
sudo journalctl -u uinetd -f

# 或直接运行查看输出
sudo uinetd -c /etc/uinetd.conf
```

### 常见问题

1. **端口已被占用**
   ```
   错误: 无法监听 TCP 0.0.0.0:80: bind: address already in use
   ```
   解决：检查端口是否被其他程序占用 `netstat -tulpn | grep :80`

2. **权限不足**
   ```
   错误: 无法监听 TCP 0.0.0.0:80: bind: permission denied
   ```
   解决：使用 sudo 运行或使用 1024 以上的端口

3. **配置文件不存在**
   ```
   错误: 无法打开配置文件
   ```
   解决：创建配置文件或使用 `-c` 参数指定正确路径

## 命令行选项

```
-c string
    配置文件路径 (默认 "/etc/uinetd.conf")
-v
    显示版本信息
```

## 开发

### 项目结构

```
uinetd/
├── main.go          # 主程序入口
├── config.go        # 配置文件解析
├── logger.go        # 日志系统
├── tcp_proxy.go     # TCP 代理实现
├── udp_proxy.go     # UDP 代理实现
├── uinetd.conf      # 示例配置文件
├── go.mod           # Go 模块文件
└── README.md        # 说明文档
```

### 编译选项

```bash
# 标准编译
go build -o uinetd

# 优化编译（减小体积）
go build -ldflags="-s -w" -o uinetd

# 交叉编译（Linux AMD64）
GOOS=linux GOARCH=amd64 go build -o uinetd-linux-amd64

# 交叉编译（Linux ARM64）
GOOS=linux GOARCH=arm64 go build -o uinetd-linux-arm64
```

## 与 rinetd 的区别

1. **协议支持**: uinetd 原生支持 UDP，rinetd 仅支持 TCP
2. **实现语言**: uinetd 使用 Go，rinetd 使用 C
3. **并发模型**: uinetd 使用 goroutine，性能更优
4. **配置文件**: 兼容 rinetd 的配置格式

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

开发者: Your Name

## 版本历史

- v1.0.0 (2025-10-07)
  - 初始版本
  - 支持 TCP/UDP 转发
  - 支持多种日志级别
  - 支持 IPv4/IPv6


