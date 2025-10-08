# uinetd - 网络重定向服务器

[![Build Status](https://github.com/evsio0n/uinetd/workflows/Build%20and%20Release/badge.svg)](https://github.com/evsio0n/uinetd/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

uinetd 是一个类似于 rinetd 的网络重定向工具，用 Go 语言编写。它可以将 TCP 和 UDP 连接从本地地址和端口重定向到另一个地址和端口。


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

## 🙏 致谢

本项目受 [rinetd](https://github.com/samhocevar/rinetd) 启发。

---

**如果这个项目对你有帮助，请给个 ⭐️ Star！**
