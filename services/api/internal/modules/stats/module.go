// Package stats 数据看板（Phase 12 交付）。
//
// 规划能力：经营数据看板、转化漏斗、留存分析。
package stats

import "github.com/yanmai/api/internal/modules"

// New 构造数据模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("stats", 12, "数据看板")
}
