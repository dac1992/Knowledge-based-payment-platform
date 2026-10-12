// Package learning 学习中心（Phase 6 交付）。
//
// 规划能力：云点播 VOD 接入（DRM）、播放器、学习进度、笔记、打卡。
package learning

import "github.com/yanmai/api/internal/modules"

// New 构造学习模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("learning", 6, "学习中心")
}
