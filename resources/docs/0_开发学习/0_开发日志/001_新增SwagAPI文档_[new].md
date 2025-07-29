# Swag API 文档

## ✅ 原因分析
- model.User 所在的包没有被 swag init 识别或扫描
- model.User 是别名或泛型（swag 不完全支持）
- 注释不规范，导致 swag 无法推断类型结构

## ✅ 解决方法
### ✅ 方法一：添加 --parseDependency 参数

```bash
swag init --parseDependency
```
这个参数会让 swag 扫描引用的外部类型（比如 model.User）所在的包。 

这个参数必须加，如果你在 response 结构体中引用了其他包的类型。

### ✅ 方法二：手动告诉 swag 哪里是 model.User
如果还是不行，可以将 model.User 所在的目录也加入 swag 的扫描范围：

```bash
swag init -g main.go --parseDependency --parseInternal
```

也可以用 --dir 参数明确告诉它要解析的路径：

```bash
swag init --parseDependency --dir ./internal/server,./internal/model
```
