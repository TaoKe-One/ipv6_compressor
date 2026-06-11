# 构建指南

## 本地开发构建

### 安装 Go

确保已安装 Go 1.21 或更高版本：

```bash
go version
```

### 安装依赖

```bash
cd ipv6-compressor-v2
go mod download
```

### 运行

```bash
go run ./cmd/gui
```

### 构建可执行文件

```bash
# macOS/Linux
go build -o ipv6-compressor ./cmd/gui

# Windows
go build -o ipv6-compressor.exe ./cmd/gui
```

## 跨平台构建

### 使用 Go 的交叉编译

```bash
# Windows 64位
GOOS=windows GOARCH=amd64 go build -o ipv6-compressor.exe ./cmd/gui

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o ipv6-compressor-amd64 ./cmd/gui

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o ipv6-compressor-arm64 ./cmd/gui

# Linux 64位
GOOS=linux GOARCH=amd64 go build -o ipv6-compressor ./cmd/gui
```

## 使用 GitHub Actions 自动构建

推送代码到 GitHub 后，GitHub Actions 会自动构建：

- Windows 版本 (每次 push)
- macOS Intel 和 Apple Silicon 版本 (每次 push)

创建 tag 后会自动创建 Release：

```bash
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
```

## 打包为独立应用（可选）

### 使用 fyne-cross 打包

`fyne-cross` 可以创建独立的 GUI 应用包。

#### 安装 fyne-cross

```bash
go install fyne.io/fyne/v2/cmd/fyne-cross@latest
```

#### 打包 Windows

```bash
fyne-cross windows -arch=amd64 -app-version=2.0.0 ./cmd/gui
```

#### 打包 macOS

```bash
fyne-cross darwin -arch=amd64,arm64 -app-version=2.0.0 ./cmd/gui
```

## 注意事项

1. **Fyne 依赖**: Fyne 需要 CGO 支持，确保系统已安装 C 编译器
2. **Windows 构建**: Windows 下构建需要 MinGW 或类似工具
3. **macOS 签名**: 分发 macOS 应用建议进行代码签名
