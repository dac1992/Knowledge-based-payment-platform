// Command server 是研脉学术知识付费平台的 API 服务入口。
//
// 启动顺序：加载配置 → 初始化日志 → 执行数据库迁移 → 建立连接池 →
// 写入种子数据（可选）→ 初始化队列 → 装配路由 → 监听并提供优雅关闭。
// @title                       研脉学术知识付费平台 API
// @version                     0.1.0
// @description                 Phase 1 工程骨架接口。所有响应统一为 `{ code, message, data, traceId }`，`code` 为 0 表示成功。
// @BasePath                    /api/v1
// @schemes                     http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 形如 `Bearer <access_token>`
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/yanmai/api/internal/config"
	"github.com/yanmai/api/internal/database"
	"github.com/yanmai/api/internal/modules/user"
	"github.com/yanmai/api/internal/router"
	"github.com/yanmai/api/pkg/logger"
	"github.com/yanmai/api/pkg/queue"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n启动失败: %v\n\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.App.Env, cfg.Log.Level)
	if err != nil {
		return err
	}
	defer func() { _ = log.Sync() }()

	log.Info("研脉学术知识付费平台 API 启动中",
		zap.String("service", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("version", config.Version),
	)

	if cfg.Database.AutoMigrate {
		if err := database.Migrate(cfg.Database.DSN, log); err != nil {
			return fmt.Errorf("数据库迁移失败: %w", err)
		}
	} else {
		log.Warn("已跳过数据库迁移（DB_AUTO_MIGRATE=false），请确认结构已由发布流程升级")
	}

	db, err := database.New(cfg, log)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := database.Close(db); cerr != nil {
			log.Warn("关闭数据库连接池失败", zap.Error(cerr))
		}
	}()

	if cfg.Database.Seed {
		if err := user.Seed(db, log); err != nil {
			return fmt.Errorf("写入种子数据失败: %w", err)
		}
	}

	qc := queue.New(queue.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, log)
	defer func() { _ = qc.Close() }()

	engine := router.New(router.Deps{
		Config: cfg,
		DB:     db,
		Logger: log,
		Queue:  qc,
	})

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("HTTP 服务已监听",
			zap.String("addr", cfg.HTTP.Addr()),
			zap.String("health", fmt.Sprintf("http://127.0.0.1:%d/health", cfg.HTTP.Port)),
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("HTTP 服务异常退出: %w", err)
	case sig := <-stop:
		log.Info("收到退出信号，开始优雅关闭", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("优雅关闭超时: %w", err)
	}

	log.Info("服务已安全退出")
	return nil
}
