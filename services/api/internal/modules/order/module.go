// Package order 订单与支付（Phase 5 交付）。
//
// 规划能力：购物车、下单、微信/支付宝支付、支付回调、权益开通、退款工单。
// 该模块是强事务区，必须依赖 PostgreSQL 的行锁与隔离级别。
package order

import "github.com/yanmai/api/internal/modules"

// New 构造订单模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("order", 5, "订单与支付")
}
