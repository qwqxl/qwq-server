# 设置为 UTF-8 编码
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# 仓库配置
$REPO_URL = "https://github.com/qwqxl/qwq-server.git"
$BRANCH = "master"

# 提交说明
$commit_msg = Read-Host "请输入提交说明（默认：Initial commit）"
if ([string]::IsNullOrWhiteSpace($commit_msg)) {
    $commit_msg = "Initial commit"
}

# 标签（可选）
$tag = Read-Host "请输入标签名称（如 v1.0.0，留空跳过打标签）"

# 初始化仓库
Write-Host "`n初始化 Git 仓库..."
git init | Out-Null

# 设置远程仓库
git remote remove origin 2>$null
git remote add origin $REPO_URL

# 添加文件并提交
git add .
git commit -m "$commit_msg"

# 切换分支
git branch -M $BRANCH

# 推送分支
git push -u origin $BRANCH

# 如果输入了标签，则添加并推送标签
if (-not [string]::IsNullOrWhiteSpace($tag)) {
    Write-Host "打标签 $tag..."
    git tag $tag
    git push origin $tag
}
else {
    Write-Host "未指定标签，跳过打标签。"
}

Write-Host "`n✅ Git 发布流程完成！"
