@echo off

REM Use current directory as working directory
set CURRENT_DIR=%~dp0
REM Remove trailing backslash
set CURRENT_DIR=%CURRENT_DIR:~0,-1%
cd /d %CURRENT_DIR%

REM Set log file path
set LOG_FILE=%CURRENT_DIR%\service_log.txt

REM Set Python server script path
set SERVER_SCRIPT=%CURRENT_DIR%\server.py

echo Current working directory: %CURRENT_DIR%
echo Log file path: %LOG_FILE%
echo Python server script: %SERVER_SCRIPT%

:pyloop
echo %date% %time% - Starting Python service... >> %LOG_FILE%
python -u "%SERVER_SCRIPT%" >> %LOG_FILE% 2>&1
echo %date% %time% - Python service stopped, restarting in 5 seconds... >> %LOG_FILE%
timeout /t 5 /nobreak > nul
goto pyloop