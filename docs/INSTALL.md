# uinetd 安装指南

## 快速安装

### 使用包管理器安装（推荐）

#### Debian/Ubuntu
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.deb
sudo dpkg -i uinetd_*_linux_amd64.deb
sudo systemctl enable --now uinetd
```

#### RHEL/CentOS/Fedora
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.rpm
sudo rpm -i uinetd_*_linux_amd64.rpm
sudo systemctl enable --now uinetd
```

#### Arch Linux
```bash
yay -S uinetd-bin
sudo systemctl enable --now uinetd
```

#### Alpine Linux
```bash
wget https://github.com/your-username/uinetd/releases/latest/download/uinetd_*_linux_amd64.apk
sudo apk add --allow-untrusted uinetd_*_linux_amd64.apk
sudo rc-update add uinetd default
sudo rc-service uinetd start
```

### 使用 Makefile 安装

```bash
# 编译并安装
sudo make install

# 安装并配置 systemd 服务
sudo make install-systemd
```

## 详细安装步骤

### 1. 编译程序

```bash
# 基本编译
make build

# 编译所有平台
make build-all

# 交叉编译特定平台
make build-linux-amd64
make build-linux-arm64
```

### 2. 安装到系统

#### 方法 A: 使用 Makefile（推荐）

```bash
# 安装程序和配置文件
sudo make install

# 这会安装:
# - /usr/local/bin/uinetd (可执行文件)
# - /etc/uinetd.conf (配置文件，如果不存在)
```

#### 方法 B: 手动安装

```bash
# 复制二进制文件
sudo cp uinetd /usr/local/bin/
sudo chmod +x /usr/local/bin/uinetd

# 复制配置文件
sudo cp uinetd.conf /etc/
```

### 3. 配置服务

编辑配置文件：

```bash
sudo vim /etc/uinetd.conf
```

添加转发规则，例如：

```
# TCP 转发
0.0.0.0      8080       192.168.1.100   80         tcp

# UDP 转发
0.0.0.0      53         8.8.8.8         53         udp

# 同时转发 TCP 和 UDP
0.0.0.0      3000       example.com     3000       all

# 设置日志级别
loglevel 4
```

### 4. 安装 systemd 服务

```bash
# 安装 systemd 服务
sudo make install-systemd

# 或手动安装
sudo cp uinetd.service /etc/systemd/system/
sudo systemctl daemon-reload
```

### 5. 启动服务

```bash
# 启用开机自启
sudo systemctl enable uinetd

# 启动服务
sudo systemctl start uinetd

# 查看状态
sudo systemctl status uinetd

# 查看日志
sudo journalctl -u uinetd -f
```

## CI/CD 部署

### GitHub Actions

项目已包含完整的 CI/CD 流程。

#### 自动构建

推送代码时自动触发构建和测试：

```bash
git push origin main
```

#### 创建发布

打 tag 触发自动发布到多个平台：

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

自动发布内容：
- ✅ GitHub Releases（二进制文件和归档）
- ✅ deb 包（Debian/Ubuntu）
- ✅ rpm 包（RHEL/CentOS/Fedora）
- ✅ apk 包（Alpine Linux）
- ✅ Arch Linux 包
- ✅ Docker 镜像
- ✅ Homebrew 包

### 使用 Ansible 部署

#### Debian/Ubuntu
创建 `deploy.yml`:

```yaml
---
- hosts: debian_servers
  become: yes
  tasks:
    - name: 下载 uinetd deb 包
      get_url:
        url: https://github.com/your-username/uinetd/releases/download/v1.0.0/uinetd_1.0.0_linux_amd64.deb
        dest: /tmp/uinetd.deb

    - name: 安装 uinetd
      apt:
        deb: /tmp/uinetd.deb
        state: present

    - name: 配置 uinetd
      template:
        src: uinetd.conf.j2
        dest: /etc/uinetd.conf
        mode: '0644'
      notify: restart uinetd

    - name: 启用并启动服务
      systemd:
        name: uinetd
        enabled: yes
        state: started

  handlers:
    - name: restart uinetd
      systemd:
        name: uinetd
        state: restarted
```

#### RHEL/CentOS/Fedora
创建 `deploy-rpm.yml`:

```yaml
---
- hosts: rhel_servers
  become: yes
  tasks:
    - name: 安装 uinetd
      yum:
        name: https://github.com/your-username/uinetd/releases/download/v1.0.0/uinetd_1.0.0_linux_amd64.rpm
        state: present

    - name: 配置 uinetd
      template:
        src: uinetd.conf.j2
        dest: /etc/uinetd.conf
        mode: '0644'
      notify: restart uinetd

    - name: 启用并启动服务
      systemd:
        name: uinetd
        enabled: yes
        state: started

  handlers:
    - name: restart uinetd
      systemd:
        name: uinetd
        state: restarted
```

部署：

```bash
ansible-playbook -i inventory deploy.yml
```

### 使用 Docker 部署

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" -o uinetd .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/uinetd /usr/local/bin/
COPY uinetd.conf /etc/
EXPOSE 80 443
CMD ["/usr/local/bin/uinetd", "-c", "/etc/uinetd.conf"]
```

构建和运行：

```bash
docker build -t uinetd:latest .
docker run -d --name uinetd \
  --network host \
  -v /etc/uinetd.conf:/etc/uinetd.conf \
  uinetd:latest
```

## 卸载

### 使用包管理器卸载

**Debian/Ubuntu:**
```bash
sudo dpkg -r uinetd
# 或完全删除包括配置文件
sudo dpkg -P uinetd
```

**RHEL/CentOS/Fedora:**
```bash
sudo rpm -e uinetd
```

**Arch Linux:**
```bash
sudo pacman -R uinetd-bin
```

**Alpine Linux:**
```bash
sudo apk del uinetd
```

### 使用 Makefile 卸载

```bash
sudo make uninstall
```

### 手动卸载

```bash
sudo systemctl stop uinetd
sudo systemctl disable uinetd
sudo rm /usr/lib/systemd/system/uinetd.service
sudo rm /usr/local/bin/uinetd
sudo systemctl daemon-reload
# 可选：删除配置文件
sudo rm /etc/uinetd.conf
```

## 更新

```bash
# 停止服务
sudo systemctl stop uinetd

# 重新编译安装
sudo make install

# 启动服务
sudo systemctl start uinetd
```

## 故障排查

### 查看服务状态

```bash
sudo systemctl status uinetd
```

### 查看实时日志

```bash
sudo journalctl -u uinetd -f
```

### 查看历史日志

```bash
sudo journalctl -u uinetd --since "1 hour ago"
```

### 测试配置文件

```bash
uinetd -c /etc/uinetd.conf
```

### 检查端口占用

```bash
sudo netstat -tulpn | grep uinetd
# 或
sudo ss -tulpn | grep uinetd
```

## 自定义安装路径

```bash
# 安装到自定义路径
sudo make install PREFIX=/opt/uinetd

# 这会安装到:
# - /opt/uinetd/bin/uinetd
# - /etc/uinetd.conf
```

## 权限要求

- 监听 1024 以下端口需要 root 权限
- systemd 服务默认以 root 用户运行
- 如需以普通用户运行，使用 1024 以上端口

## 防火墙配置

### iptables

```bash
# 允许转发端口
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -p udp --dport 53 -j ACCEPT
```

### firewalld

```bash
# 允许转发端口
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --permanent --add-port=53/udp
sudo firewall-cmd --reload
```

## SELinux 配置

如果系统启用了 SELinux：

```bash
# 临时允许
sudo setenforce 0

# 永久配置
sudo setsebool -P nis_enabled 1
```

