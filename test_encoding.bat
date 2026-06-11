@echo off
REM 测试中文编码是否正常
chcp 65001 >nul 2>&1
echo ========================================
echo 测试中文显示
echo ========================================
echo.
echo 如果您能正常看到这行中文，说明编码设置正确
echo.
ipv6_compressor.exe --help
echo.
pause
