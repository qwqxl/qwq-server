@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM 主程序路径
set "MAIN_PATH=./cmd/server/main.go"

REM 输出目录
set "OUTPUT_BASE=bin"

REM 可执行文件名称（不带后缀）
set "BINARY_NAME=qwq-server"

REM 需要打包的平台列表，格式: GOOS_GOARCH (多个用空格分开)
set "PLATFORMS=windows_amd64 linux_amd64 darwin_amd64 darwin_arm64 linux_arm64"

REM 清理输出目录
echo 清理旧目录 %OUTPUT_BASE%
if exist "%OUTPUT_BASE%" (
    rmdir /s /q "%OUTPUT_BASE%"
)
mkdir "%OUTPUT_BASE%"

for %%P in (%PLATFORMS%) do (
    REM 解析平台和架构
    for /f "tokens=1,2 delims=_" %%a in ("%%P") do (
        set "GOOS=%%a"
        set "GOARCH=%%b"

        REM 设置输出目录
        set "OUTPUT_DIR=%OUTPUT_BASE%\!GOOS!_!GOARCH!"

        echo.
        echo ================================
        echo 正在编译 !GOOS! / !GOARCH!
        echo ================================

        mkdir "!OUTPUT_DIR!"

        REM 设置文件后缀
        if "!GOOS!"=="windows" (
            set "EXE_EXT=.exe"
        ) else (
            set "EXE_EXT="
        )

        REM 禁用 CGO
        set "CGO_ENABLED=0"

        REM 编译命令
        go build -ldflags "-s -w" -trimpath -tags=jsoniter -o "!OUTPUT_DIR!\!BINARY_NAME!!EXE_EXT!" "!MAIN_PATH!"
        if errorlevel 1 (
            echo ❌ 编译失败：!GOOS!_!GOARCH!
            pause
            exit /b 1
        ) else (
            echo ✅ 编译成功：!GOOS!_!GOARCH!
        )

        REM 复制资源目录（存在则复制）
        if exist resources (
            xcopy /E /I /Y resources "!OUTPUT_DIR!\resources" >nul
        )
        if exist configs (
            xcopy /E /I /Y configs "!OUTPUT_DIR!\configs" >nul
        )
        if exist web (
            xcopy /E /I /Y web "!OUTPUT_DIR!\web" >nul
        )


        REM 复制 dockerFile 目录（如果存在且是目录）
        if exist dockerFile (
            REM 判断是不是目录，目录是存在且有子文件夹/文件
            if exist dockerFile\ (
                xcopy /E /I /Y dockerFile "!OUTPUT_DIR!\dockerFile" >nul
            )
        )

        REM 复制 Dockerfile
        echo 正在复制 Dockerfile 到 %OUTPUT_DIR% ...
        copy /Y "Dockerfile" "%OUTPUT_DIR%\Dockerfile"

    )
)

echo.
echo 🎉 所有平台打包完成！
echo 目录结构如下：
tree /f "%OUTPUT_BASE%"

endlocal
pause
