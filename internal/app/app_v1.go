package app

import (
	"context"
	"errors"
	"fmt"
	"qwqserver/pkg/authx"
	"qwqserver/pkg/httpcore"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qwqserver/internal/config"
	"qwqserver/internal/model"
	"qwqserver/pkg/cache"
	"qwqserver/pkg/database"
	"qwqserver/pkg/logger"
	"qwqserver/pkg/util/singleton"
)

type Application struct {
	Config     *config.Config
	DB         *gorm.DB
	CachePool  *cache.Pool
	HttpEngine *httpcore.Engine
	Logger     logger.Logger
	AuthX      *authx.AuthX

	initialized bool
}

var globalApps = singleton.NewSingleton[Application]()

// Initialized 获取或初始化应用实例
func Initialized(opts ...*Application) *Application {
	app, _ := globalApps.Get(func() (*Application, error) {
		return initApplication(opts...), nil
	})
	return app
}

// IsInitialized 检查应用是否已初始化
func IsInitialized() bool {
	return globalApps.IsInitialized() && Initialized().initialized
}

func initApplication(opts ...*Application) *Application {
	var app *Application
	if len(opts) > 0 && opts[0] != nil {
		app = opts[0]
	} else {
		app = &Application{}
	}

	// 设置初始化标记
	defer func() {
		app.initialized = true
	}()

	l := app.Logger
	e := app.HttpEngine

	// 初始化组件
	app.initLogger(l)
	app.initDatabase()
	app.initCache()
	app.initAuthX()
	app.initHTTPEngine(e)
	app.runMigrations()

	return app
}

func (app *Application) initLogger(l ...logger.Logger) {
	if app.Logger == nil {
		if len(l) > 0 {
			if l[0] != nil {
				app.Logger = l[0]
			}
		}
	}
}

func (app *Application) initDatabase() {
	if app.DB == nil && app.Config != nil {
		db, err := database.InitDB(&database.Config{
			Driver:              app.Config.Database.Driver,
			DSN:                 app.Config.Database.DSN,
			LogLevel:            app.Config.Database.LogLevel,
			MaxOpenConns:        app.Config.Database.MaxOpenConns,
			MaxIdleConns:        app.Config.Database.MaxIdleConns,
			ConnMaxLifetime:     app.Config.Database.ConnMaxLifetime,
			ConnMaxIdleTime:     app.Config.Database.ConnMaxIdleTime,
			HealthCheckInterval: 3 * time.Second,
		})
		if err != nil {
			app.Logger.Error("数据库初始化失败: %v", err)
		} else {
			app.DB = db
		}
	}
}

func (app *Application) initCache() {
	if app.CachePool == nil && app.Config != nil {
		app.CachePool = cache.NewPool(app.Config.Cache)
	}
}

func (app *Application) initAuthX() {
	cp := app.CachePool
	l := app.Logger
	cl, err := cp.GetClient()
	if err != nil {
		l.Warn("authX init app err: %v", err)
	}

	conf := app.Config

	// 初始化 AuthX
	ax, err := authx.New(&authx.Config{
		JWTSecret:   conf.Auth.SecretKey,
		AccessTTL:   15 * time.Minute,
		RefreshTTL:  7 * 24 * time.Hour,
		RedisPrefix: "authx:session",
		EnableSSO:   true, // 设置为 true 开启单点登录
		CacheClient: cl,
		Hooks: authx.LifecycleHooks{
			OnLogin: func(userID, platform, device string) error {
				// 在此处理登录成功后的逻辑，例如记录日志
				l.Info("login OK user_id: %v, platform: %v, device: %v", userID, platform, device)
				return nil
			},
		},
	})

	if err != nil {
		l.Warn("authX init app err: %v", err)
	}

	app.AuthX = ax
}

func (app *Application) initHTTPEngine(h ...*httpcore.Engine) {
	if app.HttpEngine == nil {
		if len(h) > 0 {
			e := h[0]
			if e != nil {
				app.HttpEngine = e
			}
		}
	}
}

func (app *Application) runMigrations() {
	if app.Config != nil && app.Config.Database.AutoMigrate && app.DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := database.AutoMigrate(
			ctx,
			&model.User{},
			&model.Post{},
			&model.Plugin{},
			// 添加其他模型...
		); err != nil {
			app.Logger.Error("数据库迁移失败: %v", err)
		}
	}
}

func (app *Application) Close() {
	// 记录关闭过程
	app.Logger.Info("开始关闭应用资源...")
	defer app.Logger.Info("应用资源关闭完成")

	// 关闭缓存连接
	if app.CachePool != nil {
		if err := app.CachePool.Close(); err != nil {
			app.Logger.Error("缓存连接关闭失败: %v", err)
		} else {
			app.Logger.Info("缓存连接已关闭")
		}
	}

	// 关闭数据库连接
	if app.DB != nil {
		if err := database.Close(); err != nil {
			app.Logger.Error("数据库连接关闭失败: %v", err)
		} else {
			app.Logger.Info("数据库连接已关闭")
		}
	}

	// 最后关闭日志系统
	if app.Logger != nil {
		app.Logger.Close()
	}
}

// Get 安全获取应用组件
func Get[T any]() (T, error) {
	return Apps[T]()
}

func GetApp() {

}

// Apps 获取应用组件（类型安全方式）
func Apps[T any]() (T, error) {
	var zero T

	if !IsInitialized() {
		return zero, errors.New("应用未初始化")
	}

	app := Initialized()

	switch any(zero).(type) {
	case *gorm.DB:
		if app.DB == nil {
			return zero, errors.New("数据库未初始化")
		}
		// 安全类型转换
		if db, ok := any(app.DB).(T); ok {
			return db, nil
		}

	case *config.Config:
		if app.Config == nil {
			return zero, errors.New("配置未初始化")
		}
		if cfg, ok := any(app.Config).(T); ok {
			return cfg, nil
		}

	case *cache.Pool:
		if app.CachePool == nil {
			return zero, errors.New("缓存池未初始化")
		}
		if pool, ok := any(app.CachePool).(T); ok {
			return pool, nil
		}

	case *cache.Client:
		if app.CachePool == nil {
			return zero, errors.New("缓存池未初始化")
		}
		cl, err := app.CachePool.GetClient()
		if err != nil {
			return zero, fmt.Errorf("获取缓存客户端失败: %w", err)
		}
		if client, ok := any(cl).(T); ok {
			return client, nil
		}

	case *authx.AuthX:
		if app.AuthX == nil {
			return zero, errors.New("AuthX未初始化")
		}
		if e, ok := any(app.AuthX).(T); ok {
			return e, nil
		}

	case *gin.Engine:
		if app.HttpEngine == nil {
			return zero, errors.New("HTTP路由未初始化")
		}
		if e, ok := any(app.HttpEngine).(T); ok {
			return e, nil
		}

	case logger.Logger:
		if app.Logger == nil {
			return zero, errors.New("日志器未初始化")
		}
		if l, ok := any(app.Logger).(T); ok {
			return l, nil
		}

	default:
		// 尝试通过类型断言获取值
		switch t := any(&zero).(type) {
		case **gorm.DB:
			if app.DB != nil {
				*t = app.DB
				return zero, nil
			}
		case **config.Config:
			if app.Config != nil {
				*t = app.Config
				return zero, nil
			}
			// 添加其他类型的处理...
		}
		return zero, fmt.Errorf("不支持的类型: %T", zero)
	}

	// 如果执行到这里，说明类型转换失败
	return zero, fmt.Errorf("类型转换失败: %T", zero)
}
