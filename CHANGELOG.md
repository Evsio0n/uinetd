# Changelog

本项目的所有重要变更都会记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- 待发布的新功能

### 变更
- 待发布的功能变更

### 修复
- 待发布的问题修复

### 移除
- 待发布的功能移除

## [1.0.0] - 2025-10-07

### 新增
- TCP 端口转发功能
- UDP 端口转发功能
- 支持同时转发 TCP 和 UDP (ALL/BOTH 协议)
- 完整的 IPv4 和 IPv6 支持
- 可配置的日志级别 (1-4)
- 自动 UDP 会话管理和清理
- systemd 服务支持
- 跨平台编译支持（Linux, macOS, Windows, FreeBSD）
- Docker 镜像支持
- GoReleaser 自动化发布
- Homebrew 安装支持
- AUR (Arch Linux) 包支持
- 完整的 CI/CD 流程
- 详细的文档和示例

### 技术细节
- 使用 Go 1.21
- 标准项目结构 (cmd, internal, configs, deploy)
- 高并发处理 (goroutine)
- 优雅的错误处理和日志记录

[Unreleased]: https://github.com/your-username/uinetd/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-username/uinetd/releases/tag/v1.0.0

