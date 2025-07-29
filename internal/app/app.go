package app

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"qwqserver/pkg/logger"
//	"qwqserver/pkg/util/singleton"
//	"time"
//
//	"gorm.io/gorm"
//
//	"qwqserver/internal/config"
//	"qwqserver/internal/model"
//	"qwqserver/pkg/cache"
//	"qwqserver/pkg/database"
//)
//
//type Application struct {
//	Config     *config.Config
//	DB         *gorm.DB
//	CachePool  *cache.Pool
//	HTTPRouter *gin.Engine
//	// tools
//	Logger             logger.Logger
//	isInitializedValue bool
//}
//
//func Get[T any]() (T, error) {
//	v, e := Apps[T]()
//	return v, e
//}
//
//func Apps[T any]() (T, error) {
//	var zero T
//	app := Initialized()
//
//	if app == nil {
//		return zero, errors.New("application 未初始化")
//	}
//
//	switch any(zero).(type) {
//	case *gorm.DB:
//		if app.DB == nil {
//			return zero, errors.New("DB 未初始化")
//		}
//		return any(app.DB).(T), nil
//
//	case *config.Config:
//		if app.Config == nil {
//			return zero, errors.New("Config 未初始化")
//		}
//		return any(app.Config).(T), nil
//
//	case *cache.Pool:
//		if app.CachePool == nil {
//			return zero, errors.New("CachePool 未初始化")
//		}
//		return any(app.CachePool).(T), nil
//
//	case *cache.Client:
//		cp, err := Apps[*cache.Pool]()
//		if err != nil {
//			return zero, err
//		}
//		cl, err := cp.GetClient()
//		if err != nil {
//			return zero, errors.New("Cache Client 未初始化")
//		}
//		return any(cl).(T), nil
//
//	case *gin.Engine:
//		if app.HTTPRouter == nil {
//			return zero, errors.New("Engine 未初始化")
//		}
//		return any(app.HTTPRouter).(T), nil
//
//	case logger.Logger:
//		if app.Logger == nil {
//			return zero, errors.New("Logger 未初始化")
//		}
//		return any(app.Logger).(T), nil
//
//	default:
//		return zero, fmt.Errorf("未注册类型: %T", zero)
//	}
//}
//
//var globalApps = singleton.NewSingleton[Application]()
//
//func Initialized(opts ...*Application) *Application {
//	var only *Application
//	var err error
//	if only, err = globalApps.Get(func() (*Application, error) {
//		return OnlyOneInit(opts...), nil
//	}); err != nil {
//		panic(err)
//	}
//
//	return only
//}
//
//func IsInitialized() bool {
//	return globalApps.IsInitialized()
//}
//
//func OnlyOneInit(apps ...*Application) *Application {
//	app := &Application{}
//	if len(apps) > 0 {
//		app = apps[0]
//	}
//
//	l := app.Logger
//
//	cfg := app.Config
//
//	// 初始化数据库
//	db, err := database.InitDB(&database.Config{
//		Driver:              cfg.Database.Driver,
//		DSN:                 cfg.Database.DSN,
//		LogLevel:            cfg.Database.LogLevel,
//		MaxOpenConns:        cfg.Database.MaxOpenConns,
//		MaxIdleConns:        cfg.Database.MaxIdleConns,
//		ConnMaxLifetime:     cfg.Database.ConnMaxLifetime,
//		ConnMaxIdleTime:     cfg.Database.ConnMaxIdleTime,
//		HealthCheckInterval: 3 * time.Second,
//	})
//	if err != nil {
//		l.Error("数据库初始化失败 Error: %v", err)
//	}
//	app.DB = db
//
//	//  初始化缓存
//	if app.CachePool == nil {
//		cachePool := cache.NewPool(cfg.Cache)
//		app.CachePool = cachePool
//	}
//
//	// 自动迁移模型
//	ctx := context.Background()
//	if cfg.Database.AutoMigrate {
//		if err := database.AutoMigrate(
//			ctx,
//			&model.User{},
//			&model.Post{},
//			&model.Plugin{},
//			// 添加其他模型...
//		); err != nil {
//			l.Error("数据库自动迁移失败 Error: %v", err)
//		}
//	}
//
//	// http
//	r := gin.New()
//	app.HTTPRouter = r
//
//	return app
//}
//
//func (app *Application) Close() {
//	// 关闭数据库连接
//	err := database.Close()
//	l := app.Logger
//	if err != nil {
//		l.Error("数据库关闭失败 Error: %v", err)
//	}
//
//	// 关闭cache连接
//	err = app.CachePool.Close()
//	if err != nil {
//		l.Error("Redis关闭失败 Error: %v", err)
//	}
//
//	// close logger conn
//	app.Logger.Close()
//}
