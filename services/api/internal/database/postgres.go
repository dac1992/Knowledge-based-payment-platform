// Package database 负责数据库连接、迁移与健康检查。
package database

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yanmai/api/internal/config"
)

// New 建立 PostgreSQL 连接并完成连通性校验。
func New(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.Database.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 newGormLogger(log, cfg),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NowFunc:                time.Now,
	})
	if err != nil {
		return nil, fmt.Errorf("连接 PostgreSQL 失败（%s）: %w", cfg.Database.SafeDSN(), err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接池失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("PostgreSQL 连通性检查失败（%s）: %w", cfg.Database.SafeDSN(), err)
	}

	log.Info("PostgreSQL 连接就绪",
		zap.String("dsn", cfg.Database.SafeDSN()),
		zap.Int("maxOpenConns", cfg.Database.MaxOpenConns),
	)
	return db, nil
}

// Ping 探测数据库可用性，供 /health 使用。
func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("数据库尚未初始化")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close 关闭连接池。
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
