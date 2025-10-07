.PHONY: all build clean install uninstall test fmt vet lint run help release snapshot

# 变量定义
BINARY_NAME=uinetd
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 安装路径
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
SYSCONFDIR=/etc
SYSTEMD_DIR=/etc/systemd/system

# 默认目标
all: build

# 编译
build:
	@echo "正在编译 $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/uinetd
	@echo "编译完成: $(BINARY_NAME)"

# 编译（带调试信息）
build-debug:
	@echo "正在编译 $(BINARY_NAME) (debug模式)..."
	go build -o $(BINARY_NAME) ./cmd/uinetd
	@echo "编译完成: $(BINARY_NAME)"

# 交叉编译
build-linux-amd64:
	@echo "正在编译 Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/uinetd

build-linux-arm64:
	@echo "正在编译 Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 ./cmd/uinetd

build-darwin-amd64:
	@echo "正在编译 macOS AMD64..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/uinetd

build-darwin-arm64:
	@echo "正在编译 macOS ARM64..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 ./cmd/uinetd

# 编译所有平台
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
	@echo "所有平台编译完成"

# 清理
clean:
	@echo "正在清理..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -rf dist/
	go clean
	@echo "清理完成"

# 安装
install: build
	@echo "正在安装 $(BINARY_NAME)..."
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	@if [ ! -f $(DESTDIR)$(SYSCONFDIR)/uinetd.conf ]; then \
		echo "正在安装配置文件..."; \
		install -d $(DESTDIR)$(SYSCONFDIR); \
		install -m 644 configs/uinetd.conf.example $(DESTDIR)$(SYSCONFDIR)/uinetd.conf; \
	else \
		echo "配置文件已存在，跳过安装"; \
	fi
	@echo "安装完成"

# 安装 systemd 服务
install-systemd: install
	@echo "正在安装 systemd 服务..."
	install -d $(DESTDIR)$(SYSTEMD_DIR)
	install -m 644 deploy/systemd/uinetd.service $(DESTDIR)$(SYSTEMD_DIR)/uinetd.service
	@if [ -z "$(DESTDIR)" ]; then \
		systemctl daemon-reload; \
		echo "systemd 服务已安装，使用以下命令启用:"; \
		echo "  sudo systemctl enable uinetd"; \
		echo "  sudo systemctl start uinetd"; \
	fi

# 卸载
uninstall:
	@echo "正在卸载 $(BINARY_NAME)..."
	@if [ -f $(DESTDIR)$(SYSTEMD_DIR)/uinetd.service ]; then \
		if [ -z "$(DESTDIR)" ]; then \
			systemctl stop uinetd 2>/dev/null || true; \
			systemctl disable uinetd 2>/dev/null || true; \
		fi; \
		rm -f $(DESTDIR)$(SYSTEMD_DIR)/uinetd.service; \
		if [ -z "$(DESTDIR)" ]; then \
			systemctl daemon-reload; \
		fi; \
	fi
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	@echo "卸载完成 (配置文件 $(SYSCONFDIR)/uinetd.conf 已保留)"

# 测试
test:
	@echo "运行测试..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 静态检查
vet:
	@echo "运行 go vet..."
	go vet ./...

# Lint (需要安装 golangci-lint)
lint:
	@echo "运行 linter..."
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint" && exit 1)
	golangci-lint run

# 运行
run: build
	./$(BINARY_NAME) -c ./configs/uinetd.conf.example

# GoReleaser 快照构建（本地测试）
snapshot:
	@echo "运行 GoReleaser 快照构建..."
	@which goreleaser > /dev/null || (echo "请先安装 goreleaser: brew install goreleaser" && exit 1)
	goreleaser build --snapshot --clean

# GoReleaser 发布（需要 GITHUB_TOKEN）
release:
	@echo "运行 GoReleaser 发布..."
	@which goreleaser > /dev/null || (echo "请先安装 goreleaser: brew install goreleaser" && exit 1)
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "错误: 请设置 GITHUB_TOKEN 环境变量"; \
		exit 1; \
	fi
	goreleaser release --clean

# 检查 GoReleaser 配置
check-release:
	@echo "检查 GoReleaser 配置..."
	@which goreleaser > /dev/null || (echo "请先安装 goreleaser: brew install goreleaser" && exit 1)
	goreleaser check

# Docker 构建
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t uinetd:latest .

# Docker 运行
docker-run: docker-build
	@echo "运行 Docker 容器..."
	docker run --rm -it \
		--network host \
		-v $(PWD)/configs/uinetd.conf.example:/etc/uinetd.conf \
		uinetd:latest

# 帮助
help:
	@echo "uinetd Makefile 使用说明"
	@echo ""
	@echo "可用目标:"
	@echo "  make build              - 编译程序"
	@echo "  make build-debug        - 编译程序 (带调试信息)"
	@echo "  make build-all          - 编译所有平台版本"
	@echo "  make build-linux-amd64  - 编译 Linux AMD64"
	@echo "  make build-linux-arm64  - 编译 Linux ARM64"
	@echo "  make clean              - 清理编译文件"
	@echo "  make install            - 安装到系统"
	@echo "  make install-systemd    - 安装程序和 systemd 服务"
	@echo "  make uninstall          - 从系统卸载"
	@echo "  make test               - 运行测试"
	@echo "  make fmt                - 格式化代码"
	@echo "  make vet                - 运行 go vet"
	@echo "  make lint               - 运行 linter"
	@echo "  make run                - 编译并运行"
	@echo "  make snapshot           - GoReleaser 快照构建"
	@echo "  make release            - GoReleaser 发布"
	@echo "  make check-release      - 检查 GoReleaser 配置"
	@echo "  make docker-build       - 构建 Docker 镜像"
	@echo "  make docker-run         - 运行 Docker 容器"
	@echo "  make help               - 显示此帮助信息"
	@echo ""
	@echo "变量:"
	@echo "  PREFIX=$(PREFIX)         - 安装前缀"
	@echo "  DESTDIR=$(DESTDIR)       - 安装目标目录"
