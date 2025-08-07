@echo off
chcp 65001 > nul
title 打印机服务监控系统

:: 设置服务器目录变量
set SERVER_DIR=D:\wwwroot\python_demo
set LOG_FILE=%SERVER_DIR%\service_log.txt

echo 正在启动打印机服务系统...
echo 服务器目录: %SERVER_DIR%
echo 日志将记录到 %LOG_FILE%

:: 启动监控脚本
start "Python服务监控" "%SERVER_DIR%\monitor_python.bat"

echo 所有服务已在单独窗口中启动
echo 关闭此窗口将终止所有服务
pause
taskkill /f /im python.exe /t