@echo off
chcp 65001 > nul

:: 使用当前目录作为工作目录
set CURRENT_DIR=%~dp0
:: 去掉路径末尾的反斜杠
set CURRENT_DIR=%CURRENT_DIR:~0,-1%
cd /d %CURRENT_DIR%

:: 设置日志文件路径
set LOG_FILE=%CURRENT_DIR%\service_log.txt

:: 设置Hi-Print程序路径（使用常规变量）
set HIPRINT_EXE=C:\Program Files\hiprint

echo 当前工作目录: %CURRENT_DIR%
echo 日志文件路径: %LOG_FILE%
echo Hi-Print程序: %HIPRINT_EXE%

:exeloop
echo %date% %time% - 正在启动Hi-Print程序... >> %LOG_FILE%
%HIPRINT_EXE% >> %LOG_FILE% 2>&1
echo %date% %time% - Hi-Print程序已停止，5秒后重启... >> %LOG_FILE%
timeout /t 5 /nobreak > nul
goto exeloop