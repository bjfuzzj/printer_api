@echo off
chcp 65001 > nul

:: 使用当前目录作为工作目录
set CURRENT_DIR=%~dp0
:: 去掉路径末尾的反斜杠
set CURRENT_DIR=%CURRENT_DIR:~0,-1%
cd /d %CURRENT_DIR%

:: 设置日志文件路径
set LOG_FILE=%CURRENT_DIR%\service_log.txt

:: 设置Python服务器脚本路径
set SERVER_SCRIPT=%CURRENT_DIR%\server.py

echo 当前工作目录: %CURRENT_DIR%
echo 日志文件路径: %LOG_FILE%
echo Python服务脚本: %SERVER_SCRIPT%

:pyloop
echo %date% %time% - 正在启动Python服务... >> %LOG_FILE%
python -u %SERVER_SCRIPT% >> %LOG_FILE% 2>&1
echo %date% %time% - Python服务已停止，5秒后重启... >> %LOG_FILE%
timeout /t 5 /nobreak > nul
goto pyloop