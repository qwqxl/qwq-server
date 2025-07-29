你想通过 **参数 Flag** 控制 `.bat` 脚本的行为，下面是一个结构清晰、支持中文、可自定义行为的 **Git 发布 `.bat` 脚本示例**：

---

## ✅ 功能说明

* 支持参数 flag：如 `-m` 指定提交信息，`-t` 指定标签，`-b` 指定分支；
* 自动添加、提交、推送代码；
* 可选添加 Git 标签；
* 全过程支持中文（UTF-8 编码，含 BOM）。

---

## 🧩 脚本示例：`deploy.bat`

```bat
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
```

---

## 🧪 示例使用方式

```bat
deploy.bat -m "修复搜索框样式问题" -t v1.0.1 -b main
```

---

## ✅ 编码建议

> 将 `.bat` 脚本保存为 **UTF-8 with BOM**，否则中文参数会乱码。

---

如你还想加：

* 自动推送到多个远程地址；
* 自动生成版本号；
* 自动构建/压缩项目；
  我也可以帮你一并加上。是否要进一步增强功能？
