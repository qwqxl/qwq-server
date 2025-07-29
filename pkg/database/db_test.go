package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"testing"
	"time"
)

func TestDBRun(t *testing.T) {
	t.Run("testDBRun", func(t *testing.T) {
		dbTest()
	})
}

func dbTest() {
	// 从你 YAML 配置转换而来：
	dsn := "qwq:123456@tcp(43.240.222.220:3306)/qwq?charset=utf8mb4&parseTime=True&loc=Local&timeout=15s"

	// 日志等级（根据你的配置：warn）
	newLogger := logger.Default.LogMode(logger.Warn)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// 是否开启预编译语句缓存
		PrepareStmt: false,
		// 禁用嵌套事务
		DisableNestedTransaction: false,
		// 命名策略（可配置表名前缀和是否用单数表名）
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // 表前缀
			SingularTable: false, // false: 多数表，true: 单数表
			NoLowerCase:   false, // false: 默认自动小写表名
		},
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 配置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取 sql.DB 对象失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)                 // 最大连接数
	sqlDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // 最大空闲时间

	// 测试连接是否成功
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping 数据库失败: %v", err)
	}

	fmt.Println("✅ 数据库连接成功")
}
