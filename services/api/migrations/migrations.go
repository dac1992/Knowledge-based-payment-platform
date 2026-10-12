// Package migrations 把 golang-migrate 的 SQL 脚本通过 go:embed 打进二进制，
// 使部署产物无需额外携带迁移文件，也避免运行时读盘的路径问题。
package migrations

import "embed"

// FS 内嵌的迁移脚本集合。
//
//go:embed *.sql
var FS embed.FS
