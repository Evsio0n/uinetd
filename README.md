# uinetd - 网络重定向服务器

[![Build Status](https://github.com/your-repo/uinetd/workflows/Build%20and%20Release/badge.svg)](https://github.com/your-repo/uinetd/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

uinetd 是一个类似于 rinetd 的网络重定向工具，用 Go 语言编写。它可以将 TCP 和 UDP 连接从本地地址和端口重定向到另一个地址和端口。

## ✨ 功能特性

- ✅ 支持 TCP 连接转发
- ✅ 支持 UDP 数据包转发
- ✅ 支持同时转发 TCP 和 UDP (ALL/BOTH 协议)
- ✅ 支持 IPv4 和 IPv6
- ✅ 可配置的日志级别
- ✅ 自动 UDP 会话管理和清理
- ✅ 高性能并发处理
- ✅ systemd 服务支持
- ✅ 跨平台编译支持

## 📁 项目结构

```
uinetd/
├── cmd/
│   └── uinetd/          # 主程序入口
│       └── main.go
├── internal/            # 内部包
│   ├── config/          # 配置解析
│   │   └── config.go
│   ├── logger/          # 日志系统
│   │   └── logger.go
│   └── proxy/           # 代理实现
│       ├── tcp_proxy.go
│       └── udp_proxy.go
├── configs/             # 配置文件
│   └── uinetd.conf.example
├── deploy/              # 部署文件
│   └── systemd/
│       └── uinetd.service
├── docs/                # 文档
│   └── INSTALL.md
├── scripts/             # 脚本
├── .github/             # GitHub Actions
│   └── workflows/
│       └── build.yml
├── Makefile            # 构建脚本
├── go.mod              # Go 模块
├── .gitignore
└── README.md
```

## 🚀 快速开始

### 安装

#### 使用包管理器安装（推荐）

**Debian/Ubuntu:**
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.deb
sudo dpkg -i uinetd_*_linux_amd64.deb
```

**RHEL/CentOS/Fedora:**
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.rpm
sudo rpm -i uinetd_*_linux_amd64.rpm
```

**Arch Linux:**
```bash
yay -S uinetd-bin
```

**Alpine Linux:**
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.apk
sudo apk add --allow-untrusted uinetd_*_linux_amd64.apk
```

**macOS (Homebrew):**
```bash
brew tap your-username/tap
brew install uinetd
```

#### 从源代码编译

```bash
# 克隆仓库
git clone https://github.com/your-username/uinetd.git
cd uinetd

# 编译
make build

# 安装到系统
sudo make install

# 安装 systemd 服务
sudo make install-systemd
```

### 配置

编辑配置文件 `/etc/uinetd.conf`:

```conf
# 转发规则格式:
# 绑定地址    绑定端口    目标地址         目标端口    协议

# TCP 转发
0.0.0.0      8080       192.168.1.100   80         tcp

# UDP 转发
0.0.0.0      53         8.8.8.8         53         udp

# 同时转发 TCP 和 UDP
0.0.0.0      3000       example.com     3000       all

# IPv6 支持
[::1]        8080       [2001:DB8::1]   80         tcp

# 日志级别 (1-4)
loglevel 4
```

### 运行

```bash
# 直接运行
sudo uinetd -c /etc/uinetd.conf

# 使用 systemd
sudo systemctl start uinetd
sudo systemctl enable uinetd

# 查看状态
sudo systemctl status uinetd

# 查看日志
sudo journalctl -u uinetd -f
```

## 📖 使用文档

### 命令行选项

```bash
uinetd [选项]

选项:
  -c string
      配置文件路径 (默认 "/etc/uinetd.conf")
  -v
      显示版本信息
```

### 支持的协议

- `tcp` - 仅 TCP 转发
- `udp` - 仅 UDP 转发
- `all` / `both` - 同时转发 TCP 和 UDP
- `raw` - 原始套接字 (暂不支持)

### 日志级别

在配置文件中设置 `loglevel`:

- **1** - 仅记录错误
- **2** - 记录错误和被禁止连接的时间
- **3** - 记录错误和被禁止连接的详细信息
- **4** - 记录所有连接的详细信息和错误（调试模式）

## 🛠️ 开发

### 编译

```bash
# 标准编译
make build

# 编译所有平台
make build-all

# 特定平台
make build-linux-amd64
make build-linux-arm64
make build-darwin-amd64
make build-darwin-arm64
```

### 测试

```bash
# 运行测试
make test

# 代码检查
make vet

# Linting (需要 golangci-lint)
make lint

# 代码格式化
make fmt
```

### 可用的 Make 目标

```bash
make help          # 显示所有可用命令
make build         # 编译程序
make clean         # 清理编译文件
make install       # 安装到系统
make uninstall     # 从系统卸载
make test          # 运行测试
make run           # 编译并运行
```

## 📦 CI/CD

项目包含 GitHub Actions 配置，自动构建和发布:

- 推送代码时自动编译和测试
- 创建 tag 时自动发布新版本
- 支持多平台交叉编译

创建发布版本:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 🔧 配置示例

### HTTP 端口转发

```conf
# 将本地 8080 端口转发到内网 Web 服务器
0.0.0.0      8080       192.168.1.100   80         tcp
```

### DNS 转发

```conf
# 转发 DNS 查询到 Google DNS
0.0.0.0      53         8.8.8.8         53         udp
0.0.0.0      53         8.8.4.4         53         udp
```

### 游戏服务器转发

```conf
# 游戏服务器通常同时使用 TCP 和 UDP
0.0.0.0      27015      game-server.com  27015     all
```

### IPv6 转发

```conf
# IPv6 到 IPv4 转发
[::1]        8080       192.168.1.100   80         tcp

# IPv6 到 IPv6 转发
[::1]        8080       [2001:DB8::1]   8080       tcp
```

## 🐛 故障排查

### 查看日志

```bash
# systemd 日志
sudo journalctl -u uinetd -f

# 指定时间范围
sudo journalctl -u uinetd --since "1 hour ago"
```

### 检查端口

```bash
# 检查端口占用
sudo netstat -tulpn | grep uinetd
sudo ss -tulpn | grep uinetd
```

### 常见问题

**端口已被占用**
```
错误: 无法监听 TCP 0.0.0.0:80: bind: address already in use
```
解决: 检查端口是否被其他程序占用 `sudo lsof -i :80`

**权限不足**
```
错误: 无法监听 TCP 0.0.0.0:80: bind: permission denied
```
解决: 使用 sudo 运行或使用 1024 以上的端口

更多信息请查看 [安装文档](docs/INSTALL.md)

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📮 联系方式

- 项目主页: https://github.com/your-repo/uinetd
- 问题反馈: https://github.com/your-repo/uinetd/issues

## 🙏 致谢

本项目受 [rinetd](https://github.com/samhocevar/rinetd) 启发。

---

**如果这个项目对你有帮助，请给个 ⭐️ Star！**
