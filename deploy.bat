@echo off
chcp 65001 >nul
setlocal ENABLEEXTENSIONS ENABLEDELAYEDEXPANSION

:: 初始化默认值
set "COMMIT_MSG=更新内容"
set "BRANCH=main"
set "TAG="

:: 参数解析
:parse
if "%~1"=="" goto after_parse
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
goto parse

:after_parse

echo.
echo ⚙️ 正在部署到 Git...
echo 📄 提交信息: %COMMIT_MSG%
echo 🌿 分支: %BRANCH%
if not "%TAG%"=="" (
    echo 🏷️ 标签: %TAG%
)

:: Git 操作开始
git init
git add .
git commit -m "%COMMIT_MSG%"
git branch -M %BRANCH%
git push -u origin %BRANCH%

if not "%TAG%"=="" (
    git tag %TAG%
    git push origin %TAG%
)

echo.
echo ✅ 推送完成！
goto :eof
