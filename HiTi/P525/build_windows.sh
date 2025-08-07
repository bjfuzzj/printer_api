#!/bin/bash
echo "正在为Windows 64位编译Go程序..."
GOOS=windows GOARCH=amd64 go build -o printer_525_status_server.exe server.go

if [ $? -eq 0 ]; then
    echo "编译成功！生成了 printer_status_server.exe"
    echo "请将以下文件复制到Windows系统："
    echo "1. printer_525_status_server.exe"
    echo "2. HTRTApiW.dll (从 Dll/Win64/ 目录)"
    echo "3. 运行 printer_525_status_server.exe 启动服务"
else
    echo "编译失败！"
fi