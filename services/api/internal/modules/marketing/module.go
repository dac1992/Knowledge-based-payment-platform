// Package marketing 营销中心（Phase 7 交付）。
//
// 规划能力：优惠券 → 拼团 → 秒杀 → 分销推广（分步交付，每个子功能独立验收）。
// 秒杀属于高并发区，依赖 PostgreSQL 行锁与 Redis 缓存。
package marketing

import "github.com/yanmai/api/internal/modules"

// New 构造营销模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("marketing", 7, "营销中心")
}
