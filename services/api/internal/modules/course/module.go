// Package course 课程内容中心（Phase 3 交付）。
//
// 规划能力：分类体系、课程 SPU/SKU、章节课时、素材上传（COS）、讲师建课、后台审核、上下架。
package course

import "github.com/yanmai/api/internal/modules"

// New 构造课程模块（骨架占位）。
func New() modules.Module {
	return modules.NewPlaceholder("course", 3, "课程内容中心")
}
