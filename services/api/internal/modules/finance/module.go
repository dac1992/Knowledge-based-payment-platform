// Package finance 财务结算（Phase 12 交付）。
//
// 规划能力：讲师分成结算、佣金结算、提现审批。
package finance

import "github.com/yanmai/api/internal/modules"

// New 构造财务模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("finance", 12, "财务结算")
}
