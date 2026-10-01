@echo off
cd /d "%~dp0backend"
echo === ACcode Go 版本測試 ===
echo.

echo 1. 檢查 Go 版本...
go version

echo.
echo 2. 清理舊的編譯檔案...
if exist accode-api.exe del accode-api.exe

echo.
echo 3. 編譯專案...
go build -o accode-api.exe .

if %errorlevel% equ 0 (
    echo.
    echo 編譯成功！
    echo.
    echo 4. 啟動伺服器...
    echo    伺服器將在 http://localhost:5001 啟動
    echo    按 Ctrl+C 停止伺服器
    echo.
    accode-api.exe
) else (
    echo.
    echo 編譯失敗！
    pause
    exit /b 1
)
