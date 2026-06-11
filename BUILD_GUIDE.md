# IPv6 Compressor Windows 打包指南

## 准备工作

您需要在 **Windows 系统** 上进行打包操作。

## 方法一：使用提供的批处理脚本（推荐）

1. 将以下文件复制到 Windows 系统：
   - `ipv6_compressor.py`
   - `build_windows.bat`
   - `requirements.txt`

2. 在 Windows 上打开命令提示符（CMD）或 PowerShell

3. 进入脚本所在目录：
   ```cmd
   cd 脚本所在目录
   ```

4. 运行打包脚本：
   ```cmd
   build_windows.bat
   ```

5. 打包完成后，可执行文件位于：`dist\ipv6_compressor.exe`

## 方法二：手动打包

如果批处理脚本无法运行，可以手动执行以下命令：

```cmd
# 1. 安装依赖
pip install -r requirements.txt

# 2. 使用 PyInstaller 打包
pyinstaller --onefile --name ipv6_compressor ipv6_compressor.py
```

## 方法三：单命令打包（最简单）

```cmd
pip install pandas tqdm openpyxl pyinstaller && pyinstaller --onefile ipv6_compressor.py
```

## 可执行文件位置

打包成功后，可执行文件位于：
```
dist\ipv6_compressor.exe
```

## 使用方法

```cmd
# 自动处理
ipv6_compressor.exe data.xlsx

# 指定输出文件
ipv6_compressor.exe data.xlsx -o result.xlsx

# 指定列名
ipv6_compressor.exe data.xlsx -c ip_address source_ip

# 使用列名模式匹配
ipv6_compressor.exe data.xlsx -p "ip|address"

# 查看帮助
ipv6_compressor.exe --help
```

## 注意事项

1. 首次运行时，可执行文件可能被 Windows Defender 误报，需要添加信任
2. 生成的 exe 文件可以独立运行，无需安装 Python
3. 文件大小约 50-100MB（包含打包的 Python 运行时）
4. 脚本已内置 Windows 中文编码处理，支持 UTF-8 输出
5. 如果仍有乱码，可在运行前设置控制台：`chcp 65001`
