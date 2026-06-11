# IPv6 Compressor v2.0

IPv6 地址批量压缩工具的 GUI 版本，将 Excel/CSV 中的 IPv6 地址转换为 RFC 5952 简写格式。

## 功能

- ✨ **友好的图形界面** - 基于 Fyne 的跨平台 GUI
- 📁 **支持多种格式** - Excel (xlsx, xls) 和 CSV 文件
- 🎯 **智能识别** - 自动检测包含 IPv6 地址的列
- 🔧 **灵活选择** - 手动选择要处理的列
- 📊 **实时进度** - 显示处理进度和统计信息
- 🖱️ **拖拽支持** - 直接拖拽文件到窗口

## 下载预编译版本

访问 [Releases](https://github.com/TaoKe-One/ipv6_compressor/releases) 页面下载：

- **Windows**: ipv6-compressor.exe
- **macOS Intel**: ipv6-compressor-amd64
- **macOS Apple Silicon**: ipv6-compressor-arm64

## 从源码构建

### 要求

- Go 1.21 或更高版本

### 构建

```bash
# 克隆仓库
git clone https://github.com/TaoKe-One/ipv6_compressor.git
cd ipv6_compressor/ipv6-compressor-v2

# 下载依赖
go mod download

# 构建
go build -o ipv6-compressor ./cmd/gui

# 运行
./ipv6-compressor
```

## 使用方法

1. **选择文件**
   - 点击"选择文件"按钮，或
   - 拖拽 Excel/CSV 文件到窗口

2. **配置选项**
   - 查看检测到的列
   - 选择要处理的列（默认选中 IPv6 列）
   - 设置输出路径（可选）

3. **开始处理**
   - 点击"开始处理"按钮
   - 等待处理完成
   - 查看处理结果

## 技术栈

- **语言**: Go 1.21+
- **GUI 框架**: [Fyne](https://fyne.io/)
- **Excel 处理**: [excelize](https://github.com/xuri/excelize)

## 项目结构

```
ipv6-compressor-v2/
├── cmd/
│   └── gui/
│       └── main.go                 # 应用入口
├── internal/
│   ├── gui/                        # GUI 实现
│   ├── ipv6/                       # IPv6 压缩核心
│   └── processor/                  # 文件处理
└── pkg/
    └── models/                     # 数据模型
```

## 许可证

MIT License
