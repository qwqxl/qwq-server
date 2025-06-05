package config

import "time"

type Database struct {
	Driver                   string        `yaml:"driver" env:"DB_DRIVER" env-default:"mysql" qwq-default:"mysql"` // 数据库驱动
	DSN                      string        `yaml:"dsn" env:"DB_DSN" qwq-default:"user:pwd@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"`
	MaxOpenConns             int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"25" qwq-default:"25"`                               // 最大打开连接数
	MaxIdleConns             int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"10" qwq-default:"10"`                               // 最大空闲连接数
	ConnMaxLifetime          time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"5m" qwq-default:"5m"`                         // 连接最大生命周期
	ConnMaxIdleTime          time.Duration `yaml:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" env-default:"1m" qwq-default:"1m"`                       // 连接最大空闲时间
	LogLevel                 string        `yaml:"log_level" env:"DB_LOG_LEVEL" env-default:"warn" qwq-default:"warn"`                                     // 日志级别
	AutoMigrate              bool          `yaml:"auto_migrate" env:"DB_AUTO_MIGRATE" env-default:"true" qwq-default:"true"`                               // 是否自动迁移
	PrepareStmt              bool          `yaml:"prepare_stmt" env:"DB_PREPARE_STMT" env-default:"false" qwq-default:"false"`                             // 是否开启预编译
	DisableNestedTransaction bool          `yaml:"disable_nested_transaction" env:"DB_DISABLE_NESTED_TRANSACTION" env-default:"false" qwq-default:"false"` // 是否禁用嵌套事务
	ConnectTimeout           time.Duration `yaml:"connect_timeout" env:"DB_CONNECT_TIMEOUT" env-default:"5s" qwq-default:"15s"`
	PingInterval             time.Duration `yaml:"ping_interval" env:"DB_PING_INTERVAL" env-default:"1m" qwq-default:"1m"`
}
