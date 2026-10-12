// Package payment 支付渠道对接（Phase 5 交付）。
//
// 规划能力：微信支付 / 支付宝渠道适配、回调验签、对账、退款。
package payment

import "github.com/yanmai/api/internal/modules"

// New 构造支付模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("payment", 5, "支付渠道对接")
}
