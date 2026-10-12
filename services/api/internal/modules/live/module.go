// Package live 直播（Phase 9 交付）。
//
// 规划能力：直播间、互动（弹幕/打赏/问答）、回放、直播带货。
package live

import "github.com/yanmai/api/internal/modules"

// New 构造直播模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("live", 9, "直播")
}
