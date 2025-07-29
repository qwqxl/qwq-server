@echo off
setlocal enabledelayedexpansion

:: 设置提交信息
set /p commitMsg=请输入提交信息（默认: auto commit）:
if "%commitMsg%"=="" set commitMsg=auto commit

:: 输出提示信息
echo.
echo ===============================
echo 开始发布流程...
echo ===============================

:: 拉取最新代码防止冲突
echo 正在拉取远程分支...
git pull origin main

if errorlevel 1 (
    echo ❌ 拉取失败，请解决冲突后再试。
    pause
    exit /b 1
)

:: 添加变更
echo 正在添加更改...
git add .

:: 提交变更
echo 正在提交...
git commit -m "%commitMsg%"

if errorlevel 1 (
    echo ⚠️ 无变更，无需提交。
) else (
    :: 推送到远程
    echo 正在推送到远程仓库...
    git push origin main

    if errorlevel 1 (
        echo ❌ 推送失败，请检查网络或权限。
        pause
        exit /b 1
    )
)

echo.
echo ✅ 发布成功！
pause
