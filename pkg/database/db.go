package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	initErr    error
	dbOnce     sync.Once
	dbMutex    sync.RWMutex // 用于保护连接池状态
	closed     bool         // 标记数据库是否已关闭
)

type Config struct {
	MaxOpenConns             int           `qwq-default:"100"`                                                                           // 最大打开连接数
	MaxIdleConns             int           `qwq-default:"10"`                                                                            // 最大空闲连接数
	ConnMaxLifetime          time.Duration `qwq-default:"30m"`                                                                           // 连接最大空闲时间
	Driver                   string        `qwq-default:"mysql"`                                                                         // 数据库驱动
	DSN                      string        `qwq-default:"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"` // 数据库连接字符串
	LogLevel                 string        `qwq-default:"warn"`                                                                          // 日志级别
	PrepareStmt              bool          `qwq-default:"true"`                                                                          // 是否使用预编译语句
	DisableNestedTransaction bool          `qwq-default:"true"`                                                                          // 是否禁用嵌套事务
	ConnMaxIdleTime          time.Duration `qwq-default:"5m"`                                                                            // 连接最大
	ConnectTimeout           time.Duration `qwq-default:"5s"`                                                                            // 连接超时
	PingInterval             time.Duration `qwq-default:"1m"`                                                                            //  Ping 间隔
	TablePrefix              string        `qwq-default:""`                                                                              // 表前缀
	SingularTable            bool          // 单数表名
	NoLowerCase              bool          // 关闭小写转换
	HealthCheckInterval      time.Duration `qwq-default:"1m"` // 健康检查间隔
	// log
	GormLogger logger.Interface //  GORM 日志接口
	LogName    string           `qwq-default:""` // 日志名
	Logger     Logger           //  自定义日志
}

type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// IsLogger 判断是否实现了Logger接口 并且不能为 nil
func IsLogger(log interface{}) bool {
	if log == nil {
		return false
	}
	_, ok := log.(Logger)
	return ok
}

// InitDB 初始化数据库连接池（线程安全）
func InitDB(cfg *Config) (*gorm.DB, error) {
	dbOnce.Do(func() {
		if err := validateConfig(cfg); err != nil {
			initErr = err
			return
		}

		maxRetries := 3
		var db *gorm.DB
		//var sqlDB *sql.DB
		var err error

		for i := 0; i <= maxRetries; i++ {
			db, _, err = connectWithRetry(cfg)
			if err == nil {
				break
			}

			if i < maxRetries {
				waitTime := time.Duration(i+1) * 2 * time.Second
				if IsLogger(cfg.Logger) {
					cfg.Logger.Error("数据库连接失败，正在重试... (尝试 %d/%d) 错误: %v", i+1, maxRetries, err)
				}
				time.Sleep(waitTime)
			}
		}

		if err != nil {
			initErr = fmt.Errorf("数据库连接失败: %w", err)
			return
		}

		dbInstance = db
		// 启动全局监控协程（确保只启动一次）
		go monitorConnection(cfg)
	})
	return dbInstance, initErr
}

// 移除context超时控制
func connectWithRetry(cfg *Config) (*gorm.DB, *sql.DB, error) {
	dialector, err := createDialector(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	gormConfig := &gorm.Config{
		Logger:                   newGormLogger(cfg.LogLevel),
		PrepareStmt:              cfg.PrepareStmt,
		DisableNestedTransaction: cfg.DisableNestedTransaction,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "",    // 表前缀
		//	SingularTable: false, // 单数表名
		//	NoLowerCase:   false, // 关闭小写转换
		//},
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("gorm open failed: %w", err)
	}

	sqlDB, err := configureConnectionPool(db, cfg)
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}

// --- 辅助函数 ---
func validateConfig(cfg *Config) error {
	if cfg.MaxOpenConns < 1 {
		return errors.New("MaxOpenConns 必须大于0")
	}
	if cfg.MaxIdleConns < 0 {
		return errors.New("MaxIdleConns 不能为负数")
	}
	if cfg.ConnMaxLifetime < time.Minute {
		return errors.New("ConnMaxLifetime 至少为1分钟")
	}
	if cfg.ConnMaxIdleTime < time.Minute {
		return errors.New("ConnMaxIdleTime 至少为1分钟")
	}
	return nil
}

func createDialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres":
		return postgres.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	default:
		return nil, fmt.Errorf("不支持的驱动: %s", driver)
	}
}

func configureConnectionPool(db *gorm.DB, cfg *Config) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// 立即测试连接
	//ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	//defer cancel()
	//if err := sqlDB.PingContext(ctx); err != nil {
	//	return nil, fmt.Errorf("连接测试失败: %w", err)
	//}

	// ⚠️ 改为无 context 限制的 Ping，避免误判连接失败
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("连接测试失败: %w", err)
	}

	return sqlDB, nil
}

func newGormLogger(level string) logger.Interface {
	gormLogLevel := getLogLevel(level)
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // （日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢查询阈值
			LogLevel:                  gormLogLevel, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略 ErrRecordNotFound 错误
			Colorful:                  true,         // 彩色打印
		},
	)
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	case "debug":
		return logger.Info // GORM的Debug级别会记录所有SQL
	default:
		return logger.Warn
	}
}

// 监控协程改为使用全局实例
func monitorConnection(cfg *Config) {
	if cfg.HealthCheckInterval <= 0 {
		if IsLogger(cfg.Logger) {
			cfg.Logger.Error("无效的健康检查间隔，跳过连接监控")
		}
		return
	}
	ticker := time.NewTicker(cfg.HealthCheckInterval)

	defer ticker.Stop()

	for range ticker.C {
		dbMutex.RLock()
		if closed {
			dbMutex.RUnlock()
			return
		}
		dbMutex.RUnlock()

		// 从全局实例获取当前连接
		db, err := GetDB()
		if err != nil {
			if IsLogger(cfg.Logger) {
				cfg.Logger.Error(cfg.LogName, "获取数据库实例失败: %v", err)
			}
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			if IsLogger(cfg.Logger) {
				cfg.Logger.Error(cfg.LogName, "获取底层连接失败: %v", err)
			}
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err != nil {
			if IsLogger(cfg.Logger) {
				cfg.Logger.Error(cfg.LogName, "数据库连接异常: %v", err)
			}
			reconnect(cfg)
		}
	}
}

func monitorConnection2(sqlDB *sql.DB, cfg *Config) {
	ticker := time.NewTicker(cfg.PingInterval)
	defer ticker.Stop()

	for range ticker.C {
		dbMutex.RLock()
		if closed {
			dbMutex.RUnlock()
			return
		}
		dbMutex.RUnlock()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := sqlDB.PingContext(ctx)
		cancel()

		if err != nil {
			if IsLogger(cfg.Logger) {
				cfg.Logger.Error(cfg.LogName, "数据库连接异常: %v", err)
			}
			// 自动重连
			reconnect(cfg)
		}
	}
}

// 重连逻辑
func reconnect(cfg *Config) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// 仅关闭旧连接池不修改全局状态
	if dbInstance != nil {
		if sqlDB, err := dbInstance.DB(); err == nil {
			sqlDB.Close() // 直接关闭连接池
		}
	}

	// 重置初始化状态
	dbOnce = sync.Once{}
	// 保持 closed=false 状态

	// 重新初始化
	_, err := InitDB(cfg)
	if err != nil && IsLogger(cfg.Logger) {
		cfg.Logger.Error(cfg.LogName, "数据库重连失败: %v", err)
	}
}

// GetDB 获取数据库实例 增加双重检查
func GetDB() (*gorm.DB, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if closed {
		return nil, errors.New("数据库连接已关闭")
	}
	if dbInstance == nil {
		return nil, errors.New("数据库未初始化")
	}
	return dbInstance, nil
}

// Close 关闭数据库连接 增加状态保护
func Close() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if closed || dbInstance == nil {
		return nil
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接失败: %w", err)
	}

	closed = true // 标记为已关闭
	return sqlDB.Close()
}

// WithTransaction 执行事务操作（支持上下文）
func WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			return fmt.Errorf("回滚失败: %w (原错误: %v)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// AutoMigrate 自动迁移模型（支持上下文）
func AutoMigrate(ctx context.Context, models ...interface{}) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	return db.WithContext(ctx).AutoMigrate(models...)
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}
