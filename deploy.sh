#!/bin/sh

# === 配置区域 ===
REPO_URL="https://github.com/qwqxl/qwq-server.git"
BRANCH="master"

# === 函数：检查是否已初始化 Git 仓库 ===
if [ ! -d ".git" ]; then
  echo "👉 当前目录未初始化 Git，正在初始化..."
  git init
else
  echo "✅ 当前已是 Git 仓库"
fi

# === 设置远程仓库 ===
git remote remove origin >/dev/null 2>&1
git remote add origin "$REPO_URL"

# === 输入提交说明 ===
printf "请输入提交说明（默认: Initial commit）："
read commit_msg
if [ -z "$commit_msg" ]; then
  commit_msg="Initial commit"
fi

# === 添加所有文件并提交 ===
git add .
git commit -m "$commit_msg"

# === 切换分支并推送 ===
git branch -M "$BRANCH"
git push -u origin "$BRANCH"

# === 可选：打标签 ===
printf "请输入标签（如 v1.0.0，留空跳过）："
read tag
if [ -n "$tag" ]; then
  git tag "$tag"
  git push origin "$tag"
  echo "📌 已打标签：$tag"
else
  echo "⚠️ 未输入标签，跳过打标签"
fi

echo
echo "🚀 部署完成！"
