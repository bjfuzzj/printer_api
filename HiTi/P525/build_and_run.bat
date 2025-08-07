@echo off
echo 正在编译Go程序...
go build -o printer_status_server.exe printer_status_server.go

if %ERRORLEVEL% EQU 0 (
    echo 编译成功！
    echo 正在启动服务器...
    printer_status_server.exe
) else (
    echo 编译失败！
    pause
)