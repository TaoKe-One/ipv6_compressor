@echo off
REM IPv6 Compressor Windows 打包脚本
REM 需要先安装依赖: pip install -r requirements.txt

echo ========================================
echo IPv6 Compressor Windows 打包工具
echo ========================================
echo.

REM 检查 PyInstaller 是否安装
python -c "import PyInstaller" 2>nul
if errorlevel 1 (
    echo [1/3] 安装 PyInstaller...
    pip install pyinstaller
) else (
    echo [1/3] PyInstaller 已安装
)

echo.
echo [2/3] 开始打包...
echo.

REM 使用 PyInstaller 打包
pyinstaller --onefile ^
    --name ipv6_compressor ^
    --hidden-import=pandas ^
    --hidden-import=openpyxl ^
    --hidden-import=tqdm ^
    ipv6_compressor.py

if errorlevel 1 (
    echo.
    echo 打包失败，尝试使用控制台模式...
    pyinstaller --onefile ^
        --name ipv6_compressor ^
        --hidden-import=pandas ^
        --hidden-import=openpyxl ^
        --hidden-import=tqdm ^
        ipv6_compressor.py
)

echo.
echo [3/3] 完成!
echo 可执行文件位于: dist\ipv6_compressor.exe
echo.

pause
