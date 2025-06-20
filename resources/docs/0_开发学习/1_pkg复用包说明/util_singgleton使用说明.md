# singleton

### singleton 使用说明

```go
package main

import (
    "fmt"
    "singleton" // 替换为你的实际包路径
)

type Database struct {
    name string
}

func main() {
    // 创建Database类型的单例容器
    dbSingleton := singleton.NewSingleton[Database]()

    // 第一次获取实例（会执行初始化）
    db1 := dbSingleton.Get(func() *Database {
        fmt.Println("Initializing database...")
        return &Database{name: "MySQL"}
    })

    // 第二次获取实例（不会重复初始化）
    db2 := dbSingleton.Get(func() *Database {
        fmt.Println("This won't execute")
        return &Database{name: "AnotherDB"}
    })

    fmt.Println(db1 == db2) // true，同一个实例
    fmt.Println(db1.name)   // "MySQL"
    fmt.Println(db2.name)   // "MySQL"
}
```