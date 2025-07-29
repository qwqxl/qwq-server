@echo off
setlocal

:: 配置区（请修改为你的仓库地址）
set "REPO_URL=https://github.com/qwqxl/qwq-server.git"
set "BRANCH=main"

:: 输入版本号作为标签（例如 v1.0.0）
set /p "TAG=请输入标签名称（如 v1.0.0，留空跳过打标签）: "

echo.
echo 初始化 Git 仓库...
git init

echo 设置远程仓库...
git remote remove origin >nul 2>&1
git remote add origin %REPO_URL%

echo 添加所有文件到暂存区...
git add .

set /p "commit_msg=请输入提交说明: "
if "%commit_msg%"=="" (
    set "commit_msg=Initial commit"
)

echo 提交到本地仓库...
git commit -m "%commit_msg%"

echo 切换到主分支...
git branch -M %BRANCH%

echo 推送到远程仓库...
git push -u origin %BRANCH%

:: 如果填写了 tag，就打标签并推送
if not "%TAG%"=="" (
    echo 打标签 %TAG%...
    git tag %TAG%
    git push origin %TAG%
) else (
    echo 未指定标签，跳过打标签。
)

echo.
echo ✅ 所有操作完成。
pause
