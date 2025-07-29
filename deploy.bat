@echo off
chcp 65001 >nul
setlocal ENABLEEXTENSIONS

:: 默认配置
set "COMMIT_MSG=更新内容"
set "BRANCH=main"
set "TAG="

:: 解析参数
:parse_args
if "%~1"=="" goto after_parse

if "%~1"=="-m" (
    shift
    set "COMMIT_MSG=%~1"
) else if "%~1"=="-t" (
    shift
    set "TAG=%~1"
) else if "%~1"=="-b" (
    shift
    set "BRANCH=%~1"
) else if "%~1"=="-h" (
    goto help
) else (
    echo ❌ 未知参数: %~1
    goto help
)
shift
goto parse_args

:after_parse

echo.
echo ⚙️ 正在部署到 Git...
echo 📄 提交说明: %COMMIT_MSG%
echo 🌿 分支名称: %BRANCH%
if not "%TAG%"=="" (
    echo 🏷️ 标签名称: %TAG%
)

:: 初始化并推送
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
goto end

:help
echo.
echo 用法: deploy.bat [-m "提交说明"] [-t v1.0.0] [-b 分支名]
echo 示例: deploy.bat -m "修复Bug" -t v1.2.3 -b main
goto end

:end
endlocal
pause
