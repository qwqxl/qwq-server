# Proto

```protobuf
// 指定使用 proto3 语法，这是 Google 推出的第 3 代 Protocol Buffers
syntax = "proto3";

// 定义当前 proto 文件属于哪个 proto 包（这个是 proto 自己的命名空间）
// 在 Go 中不太用到，但跨语言和代码结构里有意义
package hello;

// 这一行是最关键的 gRPC + Go 的 glue（胶水）！
// 它告诉 protoc：
// → 生成代码要放到 Go 项目的哪个路径
// → 生成的包名是什么（Go 里的 package 名）
//
// 拆解：
// "qwqserver/pkg/hello" 是路径（你的项目模块路径 + 子目录）
// "hello" 是生成后的 Go 文件中的 package 名
option go_package = "qwqserver/pkg/hello;hello";

// ✅ 定义一个 gRPC 服务（相当于接口定义）
// 服务名叫 Greeter，相当于 Go 里的 interface GreeterServer
service Greeter {

  // 定义服务里的一个方法：SayHello
  // 它接受一个 HelloRequest 请求
  // 返回一个 HelloReply 响应
  // 会生成 Go 接口中的方法：
  // SayHello(context.Context, *HelloRequest) (*HelloReply, error)
  rpc SayHello (HelloRequest) returns (HelloReply);
}

// ✅ 定义请求参数结构
// 相当于 Go 的 struct：
// type HelloRequest struct { Name string }
message HelloRequest {
  // 字段名是 name，类型是 string，字段编号是 1（必须唯一）
  // 这个编号用于序列化编码，不是顺序写的，是“字段标签”
  string name = 1;
}

// ✅ 定义响应结构
// 相当于 Go 的 struct：
// type HelloReply struct { Message string }
message HelloReply {
  // 字段名是 message，类型是 string，字段编号是 1
  string message = 1;
}
```


## Greeter 生成代码

```bash
protoc -I=proto \
  --go_out=pkg/hello --go_opt=paths=source_relative \
  --go-grpc_out=pkg/hello --go-grpc_opt=paths=source_relative \
  proto/hello.proto
```