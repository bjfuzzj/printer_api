@echo off

REM Use current directory as working directory
set CURRENT_DIR=%~dp0
REM Remove trailing backslash
set CURRENT_DIR=%CURRENT_DIR:~0,-1%
cd /d %CURRENT_DIR%

REM Set log file path
set LOG_FILE=%CURRENT_DIR%\service_log.txt

REM Set Hi-Print program path
set HIPRINT_EXE=C:\Program Files\hiprint\hiprint.exe

echo Current working directory: %CURRENT_DIR%
echo Log file path: %LOG_FILE%
echo Hi-Print program: %HIPRINT_EXE%
echo.
echo Starting Hi-Print monitoring...
echo Press Ctrl+C to stop
echo.

REM Check if Hi-Print executable exists
if not exist "%HIPRINT_EXE%" (
    echo ERROR: Hi-Print executable not found at: %HIPRINT_EXE%
    echo Please check the installation path
    pause
    exit /b 1
)

:exeloop
echo %date% %time% - Starting Hi-Print program...
echo %date% %time% - Starting Hi-Print program... >> %LOG_FILE%
"%HIPRINT_EXE%" >> %LOG_FILE% 2>&1
set EXIT_CODE=%errorlevel%

echo Hi-Print exited with code: %EXIT_CODE%
echo %date% %time% - Hi-Print program stopped with exit code %EXIT_CODE%, restarting in 5 seconds... >> %LOG_FILE%
echo Restarting in 5 seconds...

timeout /t 5 /nobreak > nul
goto exeloop