package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // 注册 postgres 驱动
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/yanmai/api/migrations"
)

// Migrate 执行内嵌的 SQL 迁移脚本（幂等，可重复调用）。
func Migrate(dsn string, log *zap.Logger) error {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("加载内嵌迁移脚本失败: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return fmt.Errorf("初始化迁移器失败: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Warn("关闭迁移源时出现问题", zap.Error(srcErr))
		}
		if dbErr != nil {
			log.Warn("关闭迁移数据库连接时出现问题", zap.Error(dbErr))
		}
	}()

	before, dirty, err := m.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		before = 0
	case err != nil:
		return fmt.Errorf("读取迁移版本失败: %w", err)
	}
	if dirty {
		return fmt.Errorf("迁移版本 %d 处于 dirty 状态，请人工确认后修复 schema_migrations 表", before)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("数据库结构已是最新", zap.Uint("version", before))
			return nil
		}
		return fmt.Errorf("执行迁移失败: %w", err)
	}

	after, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("读取迁移后版本失败: %w", err)
	}
	log.Info("数据库迁移完成", zap.Uint("from", before), zap.Uint("to", after))
	return nil
}
