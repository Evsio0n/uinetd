# 发布指南

本项目使用 [GoReleaser](https://goreleaser.com/) 自动化发布流程。

## 📋 前置要求

### 1. 安装 GoReleaser（本地测试用）

```bash
# macOS
brew install goreleaser

# Linux
curl -sfL https://goreleaser.com/static/run | bash
```

### 2. 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 secrets：

- `GITHUB_TOKEN` - 自动提供，用于发布到 GitHub Releases
- `HOMEBREW_TAP_GITHUB_TOKEN` - 用于发布到 Homebrew（可选）
- `AUR_KEY` - 用于发布到 AUR（可选）

## 🚀 发布流程

### 自动发布（推荐）

1. **更新版本号和更新日志**

```bash
# 确保所有更改已提交
git add .
git commit -m "chore: prepare for release vX.Y.Z"
```

2. **创建并推送 tag**

```bash
# 创建 tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送 tag 到 GitHub
git push origin v1.0.0
```

3. **自动构建和发布**

GitHub Actions 会自动：
- 运行测试
- 编译多平台二进制文件
- 构建 Docker 镜像
- 创建 GitHub Release
- 上传所有构建产物
- 生成更新日志
- 发布到 Homebrew（如已配置）
- 发布到 AUR（如已配置）

### 手动发布（本地）

```bash
# 检查配置
goreleaser check

# 测试构建（不发布）
goreleaser build --snapshot --clean

# 正式发布
export GITHUB_TOKEN="your_github_token"
goreleaser release --clean
```

## 📦 发布内容

每次发布会创建：

### 1. 二进制文件

- Linux: amd64, arm64, arm (v6, v7)
- macOS: amd64, arm64
- Windows: amd64
- FreeBSD: amd64, arm64

### 2. 归档文件

每个平台的 tar.gz/zip 文件包含：
- 可执行文件
- README.md
- LICENSE
- 配置文件示例
- systemd 服务文件
- 安装脚本

### 3. Docker 镜像

```bash
# 拉取最新版本
docker pull ghcr.io/your-username/uinetd:latest

# 拉取特定版本
docker pull ghcr.io/your-username/uinetd:v1.0.0

# 特定架构
docker pull ghcr.io/your-username/uinetd:v1.0.0-amd64
docker pull ghcr.io/your-username/uinetd:v1.0.0-arm64
```

### 4. 包管理器

**Homebrew (macOS/Linux)**
```bash
brew tap your-username/tap
brew install uinetd
```

**AUR (Arch Linux)**
```bash
yay -S uinetd-bin
```

## 🏷️ 版本命名规范

遵循 [语义化版本](https://semver.org/lang/zh-CN/)：

- `vMAJOR.MINOR.PATCH`
- 例如: `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`

**版本号说明：**
- MAJOR: 不兼容的 API 变更
- MINOR: 向下兼容的功能新增
- PATCH: 向下兼容的问题修复

**预发布版本：**
- `v1.0.0-alpha.1` - Alpha 版本
- `v1.0.0-beta.1` - Beta 版本
- `v1.0.0-rc.1` - Release Candidate

## 📝 更新日志规范

使用 [Conventional Commits](https://www.conventionalcommits.org/zh-hans/)：

```bash
# 新功能
feat: 添加 UDP 转发支持
feat(proxy): 实现 IPv6 支持

# Bug 修复
fix: 修复连接泄露问题
fix(config): 修复配置解析错误

# 性能优化
perf: 优化内存使用

# 文档
docs: 更新 README

# 测试
test: 添加单元测试

# 构建
build: 更新依赖

# CI/CD
ci: 更新 GitHub Actions
```

GoReleaser 会自动根据 commit 类型分组生成更新日志。

## 🧪 发布前检查清单

- [ ] 所有测试通过
- [ ] 代码已通过 lint 检查
- [ ] 文档已更新
- [ ] CHANGELOG.md 已更新（可选）
- [ ] 版本号符合语义化版本规范
- [ ] 本地测试构建成功

```bash
# 运行完整检查
make test
make vet
make lint

# 测试本地构建
goreleaser build --snapshot --clean
```

## 🔄 回滚发布

如果发现问题需要回滚：

1. **删除 GitHub Release**
```bash
# 在 GitHub 网页上删除 Release
# 或使用 gh CLI
gh release delete v1.0.0
```

2. **删除 tag**
```bash
# 删除本地 tag
git tag -d v1.0.0

# 删除远程 tag
git push --delete origin v1.0.0
```

3. **删除 Docker 镜像**
```bash
# 在 GitHub Packages 页面删除
# 或使用 API/CLI
```

## 📊 发布统计

查看发布统计和下载量：

- GitHub Releases: https://github.com/your-username/uinetd/releases
- Docker 镜像: https://github.com/your-username/uinetd/pkgs/container/uinetd

## 🐛 问题排查

### GoReleaser 构建失败

```bash
# 检查配置
goreleaser check

# 查看详细日志
goreleaser release --clean --debug
```

### Docker 镜像推送失败

确保已登录 GitHub Container Registry：
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### tag 已存在

```bash
# 删除并重新创建
git tag -d v1.0.0
git push --delete origin v1.0.0
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 🔗 相关链接

- [GoReleaser 文档](https://goreleaser.com/)
- [语义化版本](https://semver.org/lang/zh-CN/)
- [Conventional Commits](https://www.conventionalcommits.org/zh-hans/)
- [GitHub Actions](https://docs.github.com/en/actions)

