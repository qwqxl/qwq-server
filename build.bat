@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

set BINARY_NAME=qwq-server
set MAIN_PATH=./cmd/server/main.go

echo ===============================
echo 正在清理旧的打包目录 bin\
echo ===============================
rmdir /s /q bin 2>nul
mkdir bin

echo.
echo ===============================
echo 开始编译 Linux 版本（amd64）
echo ===============================
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -x -o bin\linux\%BINARY_NAME% -ldflags "-s -w" -trimpath -tags=jsoniter %MAIN_PATH%
if errorlevel 1 (
    echo ❌ Linux 编译失败，请检查错误信息！
    pause
    exit /b 1
)
echo ✅ Linux 编译成功！

echo.
echo 开始复制资源目录 resources 和 configs 到 Linux 打包目录
if exist resources (
    xcopy resources bin\linux\resources /E /I /Y >nul
) else (
    echo ⚠️ 找不到 resources 文件夹，跳过复制。
)
if exist configs (
    xcopy configs bin\linux\configs /E /I /Y >nul
) else (
    echo ⚠️ 找不到 configs 文件夹，跳过复制。
)

echo.
echo ===============================
echo 开始编译 Windows 版本（amd64）
echo ===============================
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -x -o bin\windows\%BINARY_NAME%.exe -ldflags "-s -w" -trimpath -tags=jsoniter %MAIN_PATH%
if errorlevel 1 (
    echo ❌ Windows 编译失败，请检查错误信息！
    pause
    exit /b 1
)
echo ✅ Windows 编译成功！

echo.
echo 开始复制资源目录 resources、configs 和 web 到 Windows 打包目录
if exist resources (
    xcopy resources bin\windows\resources /E /I /Y >nul
) else (
    echo ⚠️ 找不到 resources 文件夹，跳过复制。
)
if exist configs (
    xcopy configs bin\windows\configs /E /I /Y >nul
) else (
    echo ⚠️ 找不到 configs 文件夹，跳过复制。
)
if exist web (
    xcopy web bin\windows\web /E /I /Y >nul
) else (
    echo ⚠️ 找不到 web 文件夹，跳过复制。
)

echo.
echo ===========================
