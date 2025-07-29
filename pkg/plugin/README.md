#  plugin

## 开始

```bash
protoc -I=proto --go_out=pkg/plugin --go_opt=paths=source_relative --go-grpc_out=pkg/plugin --go-grpc_opt=paths=source_relative proto/plugin.proto
```

```bash

protoc -I=proto  --go_out=pkg/plugin/pb --go_opt=paths=source_relative  --go-grpc_out=pkg/plugin/pb --go-grpc_opt=paths=source_relative proto/plugin.proto

```


## 命令解释

| 参数	                            | 含义|
|--------------------------------| --- |
| -I=proto                       | 指定 proto 文件的根目录 |
| --go_out=pkg/plugin            | 生成 .pb.go 文件到 pkg/plugin/ |
| --go_opt=paths=source_relative | 保留 proto 原有子路径结构 |
| --go-grpc_out=...              | 生成 gRPC 服务代码到指定目录 |
| --go-grpc_opt...               | 同样保持路径结构 |
