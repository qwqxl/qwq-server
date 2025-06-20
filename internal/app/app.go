package app

import (
	"context"
	"qwqserver/internal/base"
	"qwqserver/pkg/util/singleton"
	"time"

	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
	"qwqserver/internal/config"
	"qwqserver/internal/model"
	"qwqserver/internal/server"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/database"
)

type Application struct {
	Config *config.Config
	DB     *gorm.DB
	// 其他组件...\
	Redis *redis.Client
	Log   base.Logger
	//PasswordStore *security.PasswordStore
	//PasswordSvc   security.PasswordService
}

type Config struct {
	L      base.Logger
	Config *config.Config
}

func IsLogger(l interface{}) bool {
	if l == nil {
		return false
	}
	_, ok := l.(base.Logger)
	return ok
}

var appSingleton singleton.Singleton[Application]

func New(cfg ...*Config) *Application {
	return appSingleton.Get(func() *Application {
		return NewApplication(cfg...)
	})
}

func NewApplication(appConfigs ...*Config) *Application {
	if len(appConfigs) == 0 {
		return nil
	}
	appCfg := appConfigs[0]
	l := appCfg.L
	cfg := appCfg.Config
	app := &Application{
		Config: cfg,
	}

	// 初始化日志
	//logger.Init("./logs")
	//defer logger.Close()

	// 配置数据库日志组
	//logger.ConfigureGroup("db",
	//	logger.WithAsync(false, 0),
	//	logger.WithRotateCount(5),
	//	logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
	//	logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
	//	logger.WithMaxSize(50),
	//)
	//
	//// 配置服务器日志组
	//logger.ConfigureGroup("server",
	//	logger.WithAsync(false, 0),
	//	logger.WithRotateCount(5),
	//	logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
	//	logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
	//	logger.WithMaxSize(50),
	//)
	//
	//logger.ConfigureGroup("app",
	//	logger.WithAsync(false, 0),
	//	logger.WithRotateCount(5),
	//	logger.WithFilePath("{group}/{name}/{year}/{month}/{day}"),
	//	logger.WithFilePattern("{group}_{name}_{yyyy}{mm}{dd}_{hh}{ii}{ss}_{index}.log"),
	//	logger.WithMaxSize(50),
	//)
	//
	//logger.Configure(config.DBDriverName(),
	//	logger.WithGroup("db"),
	//	logger.WithConsoleOutput(true),
	//	//logger.WithAsync(false, 0), // 关闭异步，方便测试
	//)
	//
	//logger.Configure(config.CacheDriverName(),
	//	logger.WithGroup("db"),
	//	logger.WithConsoleOutput(true),
	//)
	//
	//logger.Configure("http",
	//	logger.WithGroup("server"),
	//	logger.WithConsoleOutput(true),
	//)
	//
	//logger.Configure("auth",
	//	logger.WithGroup("server"),
	//	logger.WithConsoleOutput(true),
	//)
	//
	//logger.Configure("dev",
	//	logger.WithGroup("app"),
	//	logger.WithConsoleOutput(true),
	//)

	// 初始化数据库
	db, err := database.InitDB(&database.Config{
		Driver:              cfg.Database.Driver,
		DSN:                 cfg.Database.DSN,
		LogLevel:            cfg.Database.LogLevel,
		MaxOpenConns:        cfg.Database.MaxOpenConns,
		MaxIdleConns:        cfg.Database.MaxIdleConns,
		ConnMaxLifetime:     cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime:     cfg.Database.ConnMaxIdleTime,
		HealthCheckInterval: 3 * time.Second,
	})
	if err != nil {
		l.Error("数据库初始化失败 Error: %v", err)
	}

	//  初始化Redis缓存
	var redisClient *redis.Client
	if cfg.Redis != nil {
		//redisClient = cache.InitRedis(&cache.Config{
		//	Addr:         cfg.Redis.Addr,
		//	Password:     cfg.Redis.Password,
		//	DB:           cfg.Redis.DB,
		//	PoolSize:     cfg.Redis.PoolSize,
		//	MinIdleConns: cfg.Redis.MinIdleConns,
		//	DialTimeout:  cfg.Redis.DialTimeout,
		//	ReadTimeout:  cfg.Redis.ReadTimeout,
		//	WriteTimeout: cfg.Redis.WriteTimeout,
		//	IdleTimeout:  cfg.Redis.IdleTimeout,
		//})
		//if err != nil {
		//	l.Error("Redis初始化失败 Error: %v", err)
		//} else {
		//	app.Redis = redisClient
		//}
		redisClient = cache.InitRedis(&cache.Config{
			Addr:     cfg.Redis.Addr,
			DB:       cfg.Redis.DB,
			Password: cfg.Redis.Password,
		})
		app.Redis = redisClient
	}

	// 自动迁移模型
	ctx := context.Background()
	if cfg.Database.AutoMigrate {
		if err := database.AutoMigrate(
			ctx,
			&model.User{},
			&model.Post{},
			// 添加其他模型...
		); err != nil {
			l.Error("数据库自动迁移失败 Error: %v", err)
		}
	}

	// 初始化路由
	server.RouterApiV1()
	// 启动服务器
	//err = server.Run(cfg.ListenAddress())
	//if err != nil {
	//	l.Error("服务器启动失败 Error: %v", err)
	//	return nil
	//}
	//l.Info("服务器启动成功, address: %v", cfg.ListenAddress())

	return &Application{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
	}
}

func (app *Application) Close() {
	// 关闭数据库连接
	err := database.Close()
	if err != nil {
		app.Log.Error("数据库关闭失败 Error: %v", err)
	}
	// 关闭Redis连接
	//err = cache.CloseRedis()
	//if err != nil {
	//	app.Log.Error("Redis关闭失败 Error: %v", err)
	//}
}

func (app *Application) Listen() {
	// 启动服务器
	err := server.Run(app.Config.ListenAddress())
	if err != nil {
		app.Log.Error("服务器启动失败 Error: %v", err)
		return
	}
}
