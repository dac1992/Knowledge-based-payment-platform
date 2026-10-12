package httpx

import "github.com/yanmai/api/pkg/response"

// Envelope 是统一响应信封的类型别名，供 OpenAPI（swaggo）注解引用。
//
// 之所以在这里再导出一次，是因为 swag 解析注解时会用「文件已 import 的包名」来
// 解析类型引用；各 handler 都已 import httpx，注解写 httpx.Envelope 即可正确解析。
type Envelope = response.Body

// PageEnvelope 是分页数据容器的类型别名，用途同上。
type PageEnvelope = response.Page
