# 包管理器安装指南

uinetd 支持多种 Linux 发行版的包管理器，方便快速安装和管理。

## 📦 支持的包格式

- ✅ **deb** - Debian, Ubuntu, Linux Mint 等
- ✅ **rpm** - RHEL, CentOS, Fedora, openSUSE 等
- ✅ **apk** - Alpine Linux
- ✅ **archlinux** - Arch Linux, Manjaro 等
- ✅ **Homebrew** - macOS, Linux

## 🚀 安装方法

### Debian / Ubuntu

```bash
# 下载 deb 包
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_1.0.0_linux_amd64.deb

# 安装
sudo dpkg -i uinetd_1.0.0_linux_amd64.deb

# 启动服务
sudo systemctl enable --now uinetd
```

### RHEL / CentOS / Fedora

```bash
# 下载 rpm 包
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_1.0.0_linux_amd64.rpm

# 安装
sudo rpm -i uinetd_1.0.0_linux_amd64.rpm
# 或使用 dnf/yum
sudo dnf install uinetd_1.0.0_linux_amd64.rpm

# 启动服务
sudo systemctl enable --now uinetd
```

### Arch Linux

```bash
# 使用 AUR helper
yay -S uinetd-bin
# 或
paru -S uinetd-bin

# 启动服务
sudo systemctl enable --now uinetd
```

### Alpine Linux

```bash
# 下载 apk 包
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_1.0.0_linux_amd64.apk

# 安装
sudo apk add --allow-untrusted uinetd_1.0.0_linux_amd64.apk

# 启动服务 (OpenRC)
sudo rc-update add uinetd default
sudo rc-service uinetd start
```

### macOS (Homebrew)

```bash
# 添加 tap
brew tap your-username/tap

# 安装
brew install uinetd

# 启动服务
brew services start uinetd
```

## 📝 包安装内容

所有包都会安装以下内容：

### 二进制文件
- `/usr/bin/uinetd` (deb/rpm/apk/archlinux)
- `/usr/local/bin/uinetd` (Homebrew)

### 配置文件
- `/etc/uinetd.conf` (标记为配置文件，升级时不会覆盖)

### systemd 服务
- `/usr/lib/systemd/system/uinetd.service` (Linux)
- Homebrew 使用 launchd (macOS)

### 文档
- `/usr/share/doc/uinetd/README.md`
- `/usr/share/doc/uinetd/LICENSE` 或 `/usr/share/licenses/uinetd/LICENSE`

## 🔄 升级

### Debian / Ubuntu
```bash
# 下载新版本
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_1.1.0_linux_amd64.deb

# 升级（会自动停止服务，升级后重启）
sudo dpkg -i uinetd_1.1.0_linux_amd64.deb
```

### RHEL / CentOS / Fedora
```bash
# 下载新版本
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_1.1.0_linux_amd64.rpm

# 升级
sudo rpm -U uinetd_1.1.0_linux_amd64.rpm
# 或
sudo dnf upgrade uinetd_1.1.0_linux_amd64.rpm
```

### Arch Linux
```bash
yay -S uinetd-bin
```

### Homebrew
```bash
brew upgrade uinetd
```

## 🗑️ 卸载

### Debian / Ubuntu
```bash
# 卸载但保留配置文件
sudo dpkg -r uinetd

# 完全卸载（包括配置文件）
sudo dpkg -P uinetd
```

### RHEL / CentOS / Fedora
```bash
# 卸载
sudo rpm -e uinetd
# 或
sudo dnf remove uinetd
```

### Arch Linux
```bash
sudo pacman -R uinetd-bin
```

### Alpine Linux
```bash
sudo apk del uinetd
```

### Homebrew
```bash
brew uninstall uinetd
```

## 🔍 验证安装

安装后可以验证：

```bash
# 检查版本
uinetd -v

# 检查服务状态
sudo systemctl status uinetd

# 查看配置文件
cat /etc/uinetd.conf

# 查看日志
sudo journalctl -u uinetd -f
```

## 📋 包元数据

### 包信息
- **名称**: uinetd
- **描述**: Network redirection server - TCP & UDP port forwarding tool
- **许可证**: MIT
- **主页**: https://github.com/your-username/uinetd

### 依赖关系
- 无运行时依赖（静态编译的 Go 二进制文件）
- systemd（可选，用于服务管理）

## 🔒 安全特性

### 文件权限
- 二进制文件: `0755` (可执行)
- 配置文件: `0644` (所有人可读，仅 root 可写)
- systemd 服务: `0644`

### 服务安全
systemd 服务配置了以下安全选项：
- `NoNewPrivileges=true` - 防止权限提升
- `LimitNOFILE=65536` - 文件描述符限制
- `LimitNPROC=512` - 进程数限制

## 📊 包大小

各平台包大小参考（v1.0.0）：

| 包类型 | 架构 | 大小 |
|--------|------|------|
| deb | amd64 | ~2.5 MB |
| deb | arm64 | ~2.3 MB |
| rpm | x86_64 | ~2.5 MB |
| rpm | aarch64 | ~2.3 MB |
| apk | x86_64 | ~2.4 MB |
| archlinux | x86_64 | ~2.4 MB |

*注：实际大小可能因版本而异*

## 🐛 故障排查

### 安装失败

**错误**: `dpkg: dependency problems`
```bash
# 修复依赖
sudo apt-get install -f
```

**错误**: `rpm: warning: Header V4 signature: NOKEY`
```bash
# RPM 包已签名，可以安全安装
sudo rpm -i --nosignature uinetd_*.rpm
```

### 服务无法启动

```bash
# 检查日志
sudo journalctl -u uinetd -n 50

# 检查配置文件
sudo uinetd -c /etc/uinetd.conf

# 重新加载 systemd
sudo systemctl daemon-reload
sudo systemctl restart uinetd
```

### 配置文件位置

如果找不到配置文件：

```bash
# Debian/Ubuntu/RHEL/CentOS/Fedora/Arch/Alpine
ls -la /etc/uinetd.conf

# Homebrew (macOS)
ls -la $(brew --prefix)/etc/uinetd.conf
```

## 🔗 相关链接

- [GitHub Releases](https://github.com/your-username/uinetd/releases)
- [安装指南](INSTALL.md)
- [发布指南](RELEASE.md)
- [主文档](../README.md)

