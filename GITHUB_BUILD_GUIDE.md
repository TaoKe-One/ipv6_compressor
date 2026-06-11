# 使用 GitHub Actions 自动构建 Windows 可执行文件

## 最简单的方式（推荐）

### 步骤：

1. **将代码推送到 GitHub**

   ```bash
   # 在 Linux 上执行
   git init
   git add ipv6_compressor.py requirements.txt .github/workflows/build-windows.yml
   git commit -m "Add IPv6 compressor"
   git branch -M main
   git remote add origin https://github.com/你的用户名/ipv6-compressor.git
   git push -u origin main
   ```

2. **查看自动构建**

   - 访问你的 GitHub 仓库
   - 点击 "Actions" 标签
   - 会看到自动运行的构建任务
   - 等待几分钟后，点击任务进入详情页
   - 在 "Artifacts" 部分下载 `ipv6_compressor-windows`

3. **手动触发构建**

   - 进入 GitHub 仓库
   - 点击 "Actions" → "Build Windows Executable"
   - 点击 "Run workflow" 按钮
   - 构建完成后下载产物

---

## 方案二：使用 Wine + PyInstaller（在 Linux 上）

### 安装 Wine

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install wine wine64 python3-pip

# CentOS/RHEL
sudo yum install wine python3-pip
```

### 打包步骤

```bash
# 1. 安装 Windows 版本的 Python（通过 Wine）
wget https://www.python.org/ftp/python/3.11.0/python-3.11.0-amd64.exe
wine python-3.11.0-amd64.exe /quiet InstallAllUsers=0 PrependPath=1

# 2. 安装 PyInstaller
wine ~/.wine/drive_c/Users/$USER/AppData/Local/Programs/Python/Python311/python.exe -m pip install pyinstaller

# 3. 打包
wine ~/.wine/drive_c/Users/$USER/AppData/Local/Programs/Python/Python311/Scripts/pyinstaller.exe --onefile ipv6_compressor.py
```

> ⚠️ 此方法较为复杂，不推荐

---

## 方案三：使用虚拟机

```bash
# 下载 Windows 10 虚拟机镜像（微软提供免费开发用）
# https://developer.microsoft.com/en-us/windows/downloads/virtual-machines/

# 或使用 VirtualBox 安装 Windows
```

---

## 快速测试

如果您没有 GitHub 账号或不想推送代码，可以使用在线构建服务：

| 服务 | 网址 |
|------|------|
| Repl.it | https://replit.com |
| CodeSandbox | https://codesandbox.io |

但推荐使用 GitHub Actions，完全免费且自动化。
