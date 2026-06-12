# IPv6 Compressor v2.1.2

IPv6 地址批量处理工具的 GUI 版本，支持 Excel/CSV 文件中 IPv6 地址的压缩（RFC 5952）和扩展（完整格式）。

## 功能特性

- ✨ **友好的图形界面** - 基于 Fyne 的跨平台 GUI，完整中文支持
- 📁 **支持多种格式** - Excel (.xlsx, .xls) 和 CSV 文件
- 🎯 **智能识别** - 自动检测包含 IPv6 地址的列
- 🔧 **双模式处理** - 支持压缩（RFC 5952）和扩展（完整格式）两种模式
- 📊 **实时进度** - 显示处理进度和统计信息
- 🖱️ **拖拽支持** - 直接拖拽文件到窗口
- 🔍 **文件过滤器** - 文件选择器自动过滤支持的格式
- 📂 **目录记忆** - 记住上次使用的目录
- 💾 **自动输出路径** - 默认与原文件同目录

## 下载预编译版本

访问 [Releases](https://github.com/TaoKe-One/ipv6_compressor/releases) 页面下载：

| 平台 | 架构 | 文件名 |
|------|------|--------|
| Windows | x64 (64位) | ipv6-compressor.exe |
| Linux | x64 (64位) | ipv6-compressor-linux-amd64 |
| macOS | Intel (x64) | ipv6-compressor-macos-amd64 |
| macOS | Apple Silicon (ARM64) | ipv6-compressor-macos-arm64 |

## 使用方法

### 1. 选择文件
- 点击"选择文件"按钮
- 或直接拖拽 Excel/CSV 文件到窗口

### 2. 配置选项
- 选择处理模式：
  - **压缩 (RFC 5952)**: 将 IPv6 地址转换为简写格式
  - **扩展 (完整格式)**: 将 IPv6 地址展开为 8 组 4 位十六进制
- 查看自动检测到的列
- 设置输出路径（可选，默认与原文件同目录）

### 3. 开始处理
- 点击"开始处理"按钮
- 实时查看处理进度
- 处理完成后查看统计信息

## 从源码构建

### 系统要求

- Go 1.21 或更高版本
- CGO 支持（Fyne GUI 需要）

### Linux 依赖

```bash
sudo apt-get install libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxrender-dev
```

### 构建

```bash
# 克隆仓库
git clone https://github.com/TaoKe-One/ipv6_compressor.git
cd ipv6_compressor

# 下载依赖
go mod download

# 构建
go build -o ipv6-compressor ./cmd/gui

# 运行
./ipv6-compressor
```

### 跨平台构建

使用 `fyne-cross` 工具进行跨平台打包：

```bash
# 安装 fyne-cross
go install fyne.io/fyne/v2/cmd/fyne-cross@latest

# 构建 Windows 版本
fyne-cross windows -arch amd64 -app-id com.taokeone.ipv6compressor

# 构建 Linux 版本
fyne-cross linux -arch amd64

# 构建 macOS 版本
fyne-cross darwin -arch amd64,arm64
```

## 技术栈

- **语言**: Go 1.21+
- **GUI 框架**: [Fyne v2.7.4](https://fyne.io/)
- **Excel 处理**: [excelize v2](https://github.com/xuri/excelize)

## 项目结构

```
ipv6-compressor-v2/
├── cmd/
│   └── gui/
│       └── main.go                 # 应用入口
├── internal/
│   ├── gui/                        # GUI 实现
│   │   ├── main.go                 # 主窗口和布局
│   │   ├── components.go           # UI 组件
│   │   └── dialogs.go              # 对话框
│   ├── ipv6/                       # IPv6 处理
│   │   └── compressor.go           # 压缩/扩展核心逻辑
│   └── processor/                  # 文件处理
│       ├── excel.go                # Excel 处理
│       ├── csv.go                  # CSV 处理
│       └── detector.go             # IPv6 列检测
└── pkg/
    └── models/                     # 数据模型
```

## 更新日志

### v2.1.2 (最新)
- 🐛 修复 Preferences API 错误
- ✨ 增强文件选择过滤器
- ✨ 添加默认输出路径自动生成
- ✨ 添加目录记忆功能
- 🔧 改进错误处理和状态显示

### v2.1.1
- ✨ 新增 IPv6 扩展功能（完整格式）
- 🐛 修复 Windows 平台控制台窗口问题
- 🐛 修复文件对话框闪退问题
- 🔧 升级 Fyne 到 v2.7.4

### v2.1.0
- 🎉 初始 GUI 版本发布

## 许可证

MIT License
