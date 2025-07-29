@echo off
:: 设置终端编码为 UTF-8，防止中文乱码
chcp 65001 >nul
setlocal ENABLEEXTENSIONS ENABLEDELAYEDEXPANSION

:: 默认参数
set "COMMIT_MSG=更新"
set "BRANCH=main"
set "TAG="

:: 参数解析
:parse_args
if "%~1"=="" goto run
if /I "%~1"=="-m" (
    shift
    set "COMMIT_MSG=%~1"
) else if /I "%~1"=="-b" (
    shift
    set "BRANCH=%~1"
) else if /I "%~1"=="-t" (
    shift
    set "TAG=%~1"
) else (
    echo ❌ 未知参数: %~1
    goto :eof
)
shift
goto parse_args

:: 开始执行 Git 操作
:run
echo.
echo ========== Git 自动发布 ==========
echo 📄 提交说明：%COMMIT_MSG%
echo 🌿 分支名称：%BRANCH%
if not "%TAG%"=="" (
    echo 🏷️ 标签版本：%TAG%
)

:: 开始上传
git init
git add .
git commit -m "%COMMIT_MSG%"
git branch -M %BRANCH%
git push -u origin %BRANCH%

:: 添加并推送标签
if not "%TAG%"=="" (
    git tag %TAG%
    git push origin %TAG%
)

echo.
echo ✅ 发布完成！
goto :eof
